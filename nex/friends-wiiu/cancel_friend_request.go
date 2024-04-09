package nex_friends_wiiu

import (
	"github.com/PretendoNetwork/friends/database"
	database_wiiu "github.com/PretendoNetwork/friends/database/wiiu"
	"github.com/PretendoNetwork/friends/globals"
	notifications_wiiu "github.com/PretendoNetwork/friends/notifications/wiiu"
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	friends_wiiu "github.com/PretendoNetwork/nex-protocols-go/v2/friends-wiiu"
)

func CancelFriendRequest(err error, packet nex.PacketInterface, callID uint32, id *types.PrimitiveU64) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.FPD.InvalidArgument, "") // TODO - Add error message
	}

	connection := packet.Sender().(*nex.PRUDPConnection)

	pid, err := database_wiiu.DeleteFriendRequestAndReturnFriendPID(id.Value)
	if err != nil {
		if err == database.ErrFriendRequestNotFound {
			return nil, nex.NewError(nex.ResultCodes.FPD.InvalidMessageID, "") // TODO - Add error message
		} else {
			globals.Logger.Critical(err.Error())
			return nil, nex.NewError(nex.ResultCodes.FPD.Unknown, "") // TODO - Add error message
		}
	}

	connectedUser, ok := globals.ConnectedUsers.Get(pid)
	if ok && connectedUser != nil {
		// * This may send the friend removed notification, but they are the same.
		go notifications_wiiu.SendFriendshipRemoved(connectedUser.Connection, connection.PID())
	}

	rmcResponse := nex.NewRMCSuccess(globals.SecureEndpoint, nil)
	rmcResponse.ProtocolID = friends_wiiu.ProtocolID
	rmcResponse.MethodID = friends_wiiu.MethodCancelFriendRequest
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
