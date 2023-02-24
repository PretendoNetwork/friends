package main

import (
	"fmt"
	"os"
	"time"

	database_3ds "github.com/PretendoNetwork/friends-secure/database/3ds"
	database_wiiu "github.com/PretendoNetwork/friends-secure/database/wiiu"
	"github.com/PretendoNetwork/friends-secure/globals"
	notifications_3ds "github.com/PretendoNetwork/friends-secure/notifications/3ds"
	notifications_wiiu "github.com/PretendoNetwork/friends-secure/notifications/wiiu"
	nex "github.com/PretendoNetwork/nex-go"
)

func startNEXServer() {
	globals.NEXServer = nex.NewServer()
	globals.NEXServer.SetFragmentSize(900)
	globals.NEXServer.SetPrudpVersion(0)
	globals.NEXServer.SetKerberosKeySize(16)
	globals.NEXServer.SetKerberosPassword(os.Getenv("KERBEROS_PASSWORD"))
	globals.NEXServer.SetPingTimeout(20) // Maybe too long?
	globals.NEXServer.SetAccessKey("ridfebb9")

	globals.NEXServer.On("Data", func(packet *nex.PacketV0) {
		request := packet.RMCRequest()

		fmt.Println("==Friends - Secure==")
		fmt.Printf("Protocol ID: %#v\n", request.ProtocolID())
		fmt.Printf("Method ID: %#v\n", request.MethodID())
		fmt.Println("====================")
	})

	globals.NEXServer.On("Kick", func(packet *nex.PacketV0) {
		pid := packet.Sender().PID()

		if globals.ConnectedUsers[pid] == nil {
			return
		}

		platform := globals.ConnectedUsers[pid].Platform
		lastOnline := nex.NewDateTime(0)
		lastOnline.FromTimestamp(time.Now())

		if platform == globals.WUP {
			database_wiiu.UpdateUserLastOnlineTime(pid, lastOnline)
			notifications_wiiu.SendUserWentOfflineGlobally(packet.Sender())
		} else if platform == globals.CTR {
			database_3ds.UpdateUserLastOnlineTime(pid, lastOnline)
			notifications_3ds.SendUserWentOfflineGlobally(packet.Sender())
		}

		delete(globals.ConnectedUsers, pid)
		fmt.Println("Leaving (Kick)")
	})

	globals.NEXServer.On("Disconnect", func(packet *nex.PacketV0) {
		fmt.Println("Leaving (Disconnect)")
	})

	globals.NEXServer.On("Connect", connect)

	assignNEXProtocols()

	globals.NEXServer.Listen(":60001")
}
