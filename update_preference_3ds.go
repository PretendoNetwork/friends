package main

import (
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func updatePreference3DS(err error, client *nex.Client, callID uint32, unknown1 bool, unknown2 bool, unknown3 bool) {
	// TODO: Do something with this

	rmcResponse := nex.NewRMCResponse(nexproto.Friends3DSProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.Friends3DSMethodUpdatePreference, nil)

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
