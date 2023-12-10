package nex

import (
	"fmt"
	"os"
	"strconv"
	"time"

	database_3ds "github.com/PretendoNetwork/friends/database/3ds"
	database_wiiu "github.com/PretendoNetwork/friends/database/wiiu"
	"github.com/PretendoNetwork/friends/globals"
	notifications_3ds "github.com/PretendoNetwork/friends/notifications/3ds"
	notifications_wiiu "github.com/PretendoNetwork/friends/notifications/wiiu"
	"github.com/PretendoNetwork/friends/types"
	nex "github.com/PretendoNetwork/nex-go"
	_ "github.com/PretendoNetwork/nex-protocols-go"
)

func StartSecureServer() {
	globals.SecureServer = nex.NewPRUDPServer()
	globals.SecureServer.SecureVirtualServerPorts = []uint8{1}
	globals.SecureServer.SetFragmentSize(962)
	globals.SecureServer.SetDefaultLibraryVersion(nex.NewLibraryVersion(1, 1, 0))
	globals.SecureServer.SetKerberosPassword([]byte("password"))
	globals.SecureServer.SetKerberosKeySize(16)
	globals.SecureServer.SetAccessKey("ridfebb9")

	globals.SecureServer.OnData(func(packet nex.PacketInterface) {
		request := packet.RMCMessage()

		fmt.Println("==Friends - Secure==")
		fmt.Printf("Protocol ID: %#v\n", request.ProtocolID)
		fmt.Printf("Method ID: %#v\n", request.MethodID)
		fmt.Println("====================")
	})

	globals.SecureServer.OnClientRemoved(func(client *nex.PRUDPClient) {
		pid := client.PID().LegacyValue()

		if globals.ConnectedUsers[pid] == nil {
			return
		}

		platform := globals.ConnectedUsers[pid].Platform
		lastOnline := nex.NewDateTime(0)
		lastOnline.FromTimestamp(time.Now())

		if platform == types.WUP {
			err := database_wiiu.UpdateUserLastOnlineTime(pid, lastOnline)
			if err != nil {
				globals.Logger.Critical(err.Error())
			}

			notifications_wiiu.SendUserWentOfflineGlobally(client)
		} else if platform == types.CTR {
			err := database_3ds.UpdateUserLastOnlineTime(pid, lastOnline)
			if err != nil {
				globals.Logger.Critical(err.Error())
			}

			notifications_3ds.SendUserWentOfflineGlobally(client)
		}

		delete(globals.ConnectedUsers, pid)
		fmt.Println("Leaving (Kick)")
	})

	registerCommonSecureServerProtocols()
	registerSecureServerProtocols()

	port, _ := strconv.Atoi(os.Getenv("PN_FRIENDS_SECURE_SERVER_PORT"))
	globals.SecureServer.Listen(port)
}
