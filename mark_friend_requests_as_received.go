package main

import (
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func markFriendRequestsAsReceived(err error, client *nex.Client, callID uint32, ids []uint64) {
	for i := 0; i < len(ids); i++ {
		id := ids[i]
		setFriendRequestReceived(id)
	}

	rmcResponse := nex.NewRMCResponse(nexproto.FriendsProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.FriendsMethodMarkFriendRequestsAsReceived, nil)

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
