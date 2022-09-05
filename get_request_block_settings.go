package main

import (
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func getRequestBlockSettings(err error, client *nex.Client, callID uint32, pids []uint32) {
	settings := make([]*nexproto.PrincipalRequestBlockSetting, 0)

	// TODO:
	// Improve this. Use less database reads
	for i := 0; i < len(pids); i++ {
		requestedPID := pids[i]

		setting := nexproto.NewPrincipalRequestBlockSetting()
		setting.PID = requestedPID
		setting.IsBlocked = isFriendRequestBlocked(client.PID(), requestedPID)

		settings = append(settings, setting)
	}

	rmcResponseStream := nex.NewStreamOut(nexServer)

	rmcResponseStream.WriteListStructure(settings)

	rmcResponseBody := rmcResponseStream.Bytes()

	// Build response packet
	rmcResponse := nex.NewRMCResponse(nexproto.FriendsWiiUProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.FriendsWiiUMethodGetRequestBlockSettings, rmcResponseBody)

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewPacketV0(client, nil)

	responsePacket.SetVersion(0)
	responsePacket.SetSource(0xA1)
	responsePacket.SetDestination(0xAF)
	responsePacket.SetType(nex.DataPacket)
	responsePacket.SetPayload(rmcResponseBytes)

	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)

	nexServer.Send(responsePacket)
}
