package nex_friends_3ds

import (
	// "github.com/PretendoNetwork/friends/globals"
	nex "github.com/PretendoNetwork/nex-go"
	// friends_3ds "github.com/PretendoNetwork/nex-protocols-go/friends-3ds"
)

func RemoveFriendByLocalFriendCode(err error, packet nex.PacketInterface, callID uint32, friendLFC uint64) (*nex.RMCMessage, uint32) {
	// Respond with unimplemented, waiting for gRPC to retrieve PID from account server

	// rmcResponse := nex.NewRMCResponse(friends_3ds.ProtocolID, callID)
	// rmcResponse.SetError(nex.Errors.Core.NotImplemented)

	// rmcResponseBytes := rmcResponse.Bytes()

	// responsePacket, _ := nex.NewPRUDPPacketV0(client, nil)

	//
	// responsePacket.SetSource(0xA1)
	// responsePacket.SetDestination(0xAF)
	// responsePacket.SetType(nex.DataPacket)
	// responsePacket.SetPayload(rmcResponseBytes)

	// responsePacket.AddFlag(nex.FlagNeedsAck)
	// responsePacket.AddFlag(nex.FlagReliable)

	// globals.SecureServer.Send(responsePacket)

	return nil, nex.Errors.Core.NotImplemented
}
