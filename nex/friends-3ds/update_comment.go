package nex_friends_3ds

import (
	database_3ds "github.com/PretendoNetwork/friends/database/3ds"
	"github.com/PretendoNetwork/friends/globals"
	notifications_3ds "github.com/PretendoNetwork/friends/notifications/3ds"
	nex "github.com/PretendoNetwork/nex-go"
	friends_3ds "github.com/PretendoNetwork/nex-protocols-go/friends-3ds"
)

func UpdateComment(err error, client *nex.Client, callID uint32, comment string) uint32 {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nex.Errors.FPD.InvalidArgument
	}

	err = database_3ds.UpdateUserComment(client.PID(), comment)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return nex.Errors.FPD.Unknown
	}

	go notifications_3ds.SendCommentUpdate(client, comment)


	rmcResponse := nex.NewRMCResponse(friends_3ds.ProtocolID, callID)
	rmcResponse.SetSuccess(friends_3ds.MethodUpdateComment, nil)

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
