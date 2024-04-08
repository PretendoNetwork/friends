package nex_friends_3ds

import (
	// "github.com/PretendoNetwork/friends/globals"
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	// friends_3ds "github.com/PretendoNetwork/nex-protocols-go/v2/friends-3ds"
)

func GetPrincipalIDByLocalFriendCode(err error, packet nex.PacketInterface, callID uint32, lfc *types.PrimitiveU64, lfcList *types.List[*types.PrimitiveU64]) (*nex.RMCMessage, *nex.Error) {
	// Respond with unimplemented, waiting for gRPC to retrieve PID from account server

	// rmcResponse := nex.NewRMCResponse(friends_3ds.ProtocolID, callID)
	// rmcResponse.SetError(nex.ResultCodes.Core.NotImplemented)

	// rmcResponseBytes := rmcResponse.Bytes()

	// responsePacket, _ := nex.NewPRUDPPacketV0(connection, nil)

	//
	// responsePacket.SetSource(0xA1)
	// responsePacket.SetDestination(0xAF)
	// responsePacket.SetType(nex.DataPacket)
	// responsePacket.SetPayload(rmcResponseBytes)

	// responsePacket.AddFlag(nex.FlagNeedsAck)
	// responsePacket.AddFlag(nex.FlagReliable)

	// globals.SecureServer.Send(responsePacket)

	return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "") // TODO - Add error message
}
