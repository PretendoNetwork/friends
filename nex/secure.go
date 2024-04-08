package nex

import (
	"os"
	"strconv"
	"time"

	database_3ds "github.com/PretendoNetwork/friends/database/3ds"
	database_wiiu "github.com/PretendoNetwork/friends/database/wiiu"
	"github.com/PretendoNetwork/friends/globals"
	notifications_3ds "github.com/PretendoNetwork/friends/notifications/3ds"
	notifications_wiiu "github.com/PretendoNetwork/friends/notifications/wiiu"
	friends_types "github.com/PretendoNetwork/friends/types"
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	_ "github.com/PretendoNetwork/nex-protocols-go"
)

func StartSecureServer() {
	port, _ := strconv.Atoi(os.Getenv("PN_FRIENDS_SECURE_SERVER_PORT"))

	globals.SecureServer = nex.NewPRUDPServer()
	globals.SecureEndpoint = nex.NewPRUDPEndPoint(1)

	globals.SecureEndpoint.AccountDetailsByPID = globals.AccountDetailsByPID
	globals.SecureEndpoint.AccountDetailsByUsername = globals.AccountDetailsByUsername
	globals.SecureEndpoint.ServerAccount = nex.NewAccount(types.NewPID(2), "Quazal Rendez-Vous", os.Getenv("PN_FRIENDS_CONFIG_SECURE_PASSWORD"))

	globals.GuestAccount = nex.NewAccount(types.NewPID(100), "guest", "MMQea3n!fsik") // * Guest account password is always the same, known to all consoles

	globals.SecureEndpoint.OnConnectionEnded(func(connection *nex.PRUDPConnection) {
		pid := connection.PID().LegacyValue()

		if globals.ConnectedUsers[pid] == nil {
			return
		}

		platform := globals.ConnectedUsers[pid].Platform
		lastOnline := types.NewDateTime(0)
		lastOnline.FromTimestamp(time.Now())

		if platform == friends_types.WUP {
			err := database_wiiu.UpdateUserLastOnlineTime(pid, lastOnline)
			if err != nil {
				globals.Logger.Critical(err.Error())
			}

			notifications_wiiu.SendUserWentOfflineGlobally(connection)
		} else if platform == friends_types.CTR {
			err := database_3ds.UpdateUserLastOnlineTime(pid, lastOnline)
			if err != nil {
				globals.Logger.Critical(err.Error())
			}

			notifications_3ds.SendUserWentOfflineGlobally(connection)
		}

		delete(globals.ConnectedUsers, pid)
	})

	registerCommonSecureServerProtocols()
	registerSecureServerProtocols()

	globals.SecureEndpoint.IsSecureEndPoint = true
	globals.SecureServer.SetFragmentSize(962)
	globals.SecureServer.LibraryVersions.SetDefault(nex.NewLibraryVersion(1, 1, 0))
	globals.SecureServer.SessionKeyLength = 16
	globals.SecureServer.AccessKey = "ridfebb9"
	globals.SecureServer.BindPRUDPEndPoint(globals.SecureEndpoint)
	globals.SecureServer.Listen(port)
}
