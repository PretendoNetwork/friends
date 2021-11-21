package main

import (
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func updateCommentWiiU(err error, client *nex.Client, callID uint32, comment *nexproto.Comment) {
	// TODO: Do something with this

	changed := updateUserComment(client.PID(), comment.Contents)

	rmcResponseStream := nex.NewStreamOut(nexServer)

	rmcResponseStream.WriteUInt64LE(changed)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(nexproto.FriendsProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.FriendsMethodUpdateComment, rmcResponseBody)

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
