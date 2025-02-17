package nex_friends_3ds

import (
	"github.com/PretendoNetwork/nex-go/v2/types"

	database_3ds "github.com/PretendoNetwork/friends/database/3ds"
	"github.com/PretendoNetwork/friends/globals"
	notifications_3ds "github.com/PretendoNetwork/friends/notifications/3ds"
	nex "github.com/PretendoNetwork/nex-go/v2"
	friends_3ds "github.com/PretendoNetwork/nex-protocols-go/v2/friends-3ds"
)

func UpdateComment(err error, packet nex.PacketInterface, callID uint32, comment types.String) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.FPD.InvalidArgument, "") // TODO - Add error message
	}

	connection := packet.Sender().(*nex.PRUDPConnection)

	err = database_3ds.UpdateUserComment(uint32(connection.PID()), string(comment))
	if err != nil {
		globals.Logger.Critical(err.Error())
		return nil, nex.NewError(nex.ResultCodes.FPD.Unknown, "") // TODO - Add error message
	}

	go notifications_3ds.SendCommentUpdate(connection, string(comment))

	rmcResponse := nex.NewRMCSuccess(globals.SecureEndpoint, nil)
	rmcResponse.ProtocolID = friends_3ds.ProtocolID
	rmcResponse.MethodID = friends_3ds.MethodUpdateComment
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
