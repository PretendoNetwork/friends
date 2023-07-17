package nex

import (
	"fmt"
	"os"
	"time"

	database_3ds "github.com/PretendoNetwork/friends-secure/database/3ds"
	database_wiiu "github.com/PretendoNetwork/friends-secure/database/wiiu"
	"github.com/PretendoNetwork/friends-secure/globals"
	notifications_3ds "github.com/PretendoNetwork/friends-secure/notifications/3ds"
	notifications_wiiu "github.com/PretendoNetwork/friends-secure/notifications/wiiu"
	"github.com/PretendoNetwork/friends-secure/types"
	nex "github.com/PretendoNetwork/nex-go"
	_ "github.com/PretendoNetwork/nex-protocols-go"
)

func StartSecureServer() {
	globals.SecureServer = nex.NewServer()
	globals.SecureServer.SetFragmentSize(900)
	globals.SecureServer.SetPRUDPVersion(0)
	globals.SecureServer.SetKerberosKeySize(16)
	globals.SecureServer.SetKerberosPassword(globals.KerberosPassword)
	globals.SecureServer.SetPingTimeout(20) // Maybe too long?
	globals.SecureServer.SetAccessKey("ridfebb9")
	globals.AuthenticationServer.SetPRUDPProtocolMinorVersion(0) // TODO: Figure out what to put here
	globals.AuthenticationServer.SetDefaultNEXVersion(&nex.NEXVersion{
		Major: 1,
		Minor: 1,
		Patch: 0,
	})

	globals.SecureServer.On("Data", func(packet *nex.PacketV0) {
		request := packet.RMCRequest()

		fmt.Println("==Friends - Secure==")
		fmt.Printf("Protocol ID: %#v\n", request.ProtocolID())
		fmt.Printf("Method ID: %#v\n", request.MethodID())
		fmt.Println("====================")
	})

	globals.SecureServer.On("Kick", func(packet *nex.PacketV0) {
		pid := packet.Sender().PID()

		if globals.ConnectedUsers[pid] == nil {
			return
		}

		platform := globals.ConnectedUsers[pid].Platform
		lastOnline := nex.NewDateTime(0)
		lastOnline.FromTimestamp(time.Now())

		if platform == types.WUP {
			database_wiiu.UpdateUserLastOnlineTime(pid, lastOnline)
			notifications_wiiu.SendUserWentOfflineGlobally(packet.Sender())
		} else if platform == types.CTR {
			database_3ds.UpdateUserLastOnlineTime(pid, lastOnline)
			notifications_3ds.SendUserWentOfflineGlobally(packet.Sender())
		}

		delete(globals.ConnectedUsers, pid)
		fmt.Println("Leaving (Kick)")
	})

	globals.SecureServer.On("Disconnect", func(packet *nex.PacketV0) {
		fmt.Println("Leaving (Disconnect)")
	})

	globals.SecureServer.On("Connect", connect)

	registerNEXProtocols()

	globals.SecureServer.Listen(":" + os.Getenv("PN_FRIENDS_SECURE_SERVER_PORT"))
}
