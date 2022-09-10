package main

import (
	"fmt"
	"os"
	"time"

	"github.com/PretendoNetwork/friends-secure/database"
	"github.com/PretendoNetwork/friends-secure/globals"
	nex "github.com/PretendoNetwork/nex-go"
)

var nexServer *nex.Server

func nexServerStart() {
	nexServer = nex.NewServer()
	nexServer.SetFragmentSize(900)
	nexServer.SetPrudpVersion(0)
	nexServer.SetKerberosKeySize(16)
	nexServer.SetKerberosPassword(os.Getenv("KERBEROS_PASSWORD"))
	nexServer.SetPingTimeout(20) // Maybe too long?
	nexServer.SetAccessKey("ridfebb9")

	nexServer.On("Data", func(packet *nex.PacketV0) {
		request := packet.RMCRequest()

		fmt.Println("==Friends - Secure==")
		fmt.Printf("Protocol ID: %#v\n", request.ProtocolID())
		fmt.Printf("Method ID: %#v\n", request.MethodID())
		fmt.Println("====================")
	})

	nexServer.On("Kick", func(packet *nex.PacketV0) {
		pid := packet.Sender().PID()
		delete(globals.ConnectedUsers, pid)

		lastOnline := nex.NewDateTime(0)
		lastOnline.FromTimestamp(time.Now())

		database.UpdateUserLastOnlineTime(pid, lastOnline)
		sendUserWentOfflineWiiUNotifications(packet.Sender())

		fmt.Println("Leaving")
	})

	nexServer.On("Ping", func(packet *nex.PacketV0) {
		fmt.Print("Pinged. Is ACK: ")
		fmt.Println(packet.HasFlag(nex.FlagAck))
	})

	nexServer.On("Connect", connect)

	assignNEXProtocols()

	nexServer.Listen(":60001")
}
