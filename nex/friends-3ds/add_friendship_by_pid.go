package nex_friends_3ds

import (
	database_3ds "github.com/PretendoNetwork/friends/database/3ds"
	"github.com/PretendoNetwork/friends/globals"
	notifications_3ds "github.com/PretendoNetwork/friends/notifications/3ds"
	nex "github.com/PretendoNetwork/nex-go"
	friends_3ds "github.com/PretendoNetwork/nex-protocols-go/friends-3ds"
)

func AddFriendshipByPrincipalID(err error, client *nex.Client, callID uint32, lfc uint64, pid uint32) uint32 {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nex.Errors.FPD.InvalidArgument
	}

	friendRelationship, err := database_3ds.SaveFriendship(client.PID(), pid)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return nex.Errors.FPD.Unknown
	}

	connectedUser := globals.ConnectedUsers[pid]
	if connectedUser != nil {
		go notifications_3ds.SendFriendshipCompleted(connectedUser.Client, pid, client.PID())
	}

	rmcResponseStream := nex.NewStreamOut(globals.SecureServer)

	rmcResponseStream.WriteStructure(friendRelationship)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(friends_3ds.ProtocolID, callID)
	rmcResponse.SetSuccess(friends_3ds.MethodAddFriendByPrincipalID, rmcResponseBody)

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewPacketV0(client, nil)

	responsePacket.SetVersion(0)
	responsePacket.SetSource(0xA1)
	responsePacket.SetDestination(0xAF)
	responsePacket.SetType(nex.DataPacket)
	responsePacket.SetPayload(rmcResponseBytes)

	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)

	globals.SecureServer.Send(responsePacket)

	return 0
}
