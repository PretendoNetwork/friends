package nex_friends_wiiu

import (
	database_wiiu "github.com/PretendoNetwork/friends-secure/database/wiiu"
	"github.com/PretendoNetwork/friends-secure/globals"
	nex "github.com/PretendoNetwork/nex-go"
	friends_wiiu "github.com/PretendoNetwork/nex-protocols-go/friends/wiiu"
)

func GetRequestBlockSettings(err error, client *nex.Client, callID uint32, pids []uint32) {
	settings := make([]*friends_wiiu.PrincipalRequestBlockSetting, 0)

	// TODO:
	// Improve this. Use less database_wiiu.reads
	for i := 0; i < len(pids); i++ {
		requestedPID := pids[i]

		setting := friends_wiiu.NewPrincipalRequestBlockSetting()
		setting.PID = requestedPID
		setting.IsBlocked = database_wiiu.IsFriendRequestBlocked(client.PID(), requestedPID)

		settings = append(settings, setting)
	}

	rmcResponseStream := nex.NewStreamOut(globals.NEXServer)

	rmcResponseStream.WriteListStructure(settings)

	rmcResponseBody := rmcResponseStream.Bytes()

	// Build response packet
	rmcResponse := nex.NewRMCResponse(friends_wiiu.ProtocolID, callID)
	rmcResponse.SetSuccess(friends_wiiu.MethodGetRequestBlockSettings, rmcResponseBody)

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewPacketV0(client, nil)

	responsePacket.SetVersion(0)
	responsePacket.SetSource(0xA1)
	responsePacket.SetDestination(0xAF)
	responsePacket.SetType(nex.DataPacket)
	responsePacket.SetPayload(rmcResponseBytes)

	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)

	globals.NEXServer.Send(responsePacket)
}
