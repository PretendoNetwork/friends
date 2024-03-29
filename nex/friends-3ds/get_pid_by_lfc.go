package nex_friends_3ds

import (
	// "github.com/PretendoNetwork/friends/globals"
	nex "github.com/PretendoNetwork/nex-go"
	// friends_3ds "github.com/PretendoNetwork/nex-protocols-go/friends-3ds"
)

func GetPrincipalIDByLocalFriendCode(err error, client *nex.Client, callID uint32, lfc uint64, lfcList []uint64) uint32 {
	// Respond with unimplemented, waiting for gRPC to retrieve PID from account server

	// rmcResponse := nex.NewRMCResponse(friends_3ds.ProtocolID, callID)
	// rmcResponse.SetError(nex.Errors.Core.NotImplemented)

	// rmcResponseBytes := rmcResponse.Bytes()

	// responsePacket, _ := nex.NewPacketV0(client, nil)

	// responsePacket.SetVersion(0)
	// responsePacket.SetSource(0xA1)
	// responsePacket.SetDestination(0xAF)
	// responsePacket.SetType(nex.DataPacket)
	// responsePacket.SetPayload(rmcResponseBytes)

	// responsePacket.AddFlag(nex.FlagNeedsAck)
	// responsePacket.AddFlag(nex.FlagReliable)

	// globals.SecureServer.Send(responsePacket)

	return nex.Errors.Core.NotImplemented
}
