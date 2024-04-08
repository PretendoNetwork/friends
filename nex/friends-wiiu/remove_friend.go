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

func RemoveFriend(err error, packet nex.PacketInterface, callID uint32, pid *types.PID) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.FPD.InvalidArgument, "") // TODO - Add error message
	}

	connection := packet.Sender().(*nex.PRUDPConnection)

	err = database_wiiu.RemoveFriendship(connection.PID().LegacyValue(), pid.LegacyValue())
	if err != nil {
		if err == database.ErrFriendshipNotFound {
			return nil, nex.NewError(nex.ResultCodes.FPD.NotInMyFriendList, "") // TODO - Add error message
		} else {
			globals.Logger.Critical(err.Error())
			return nil, nex.NewError(nex.ResultCodes.FPD.Unknown, "") // TODO - Add error message
		}
	}

	connectedUser := globals.ConnectedUsers[pid.LegacyValue()]
	if connectedUser != nil {
		go notifications_wiiu.SendFriendshipRemoved(connectedUser.Connection, pid)
	}

	rmcResponse := nex.NewRMCSuccess(globals.SecureEndpoint, nil)
	rmcResponse.ProtocolID = friends_wiiu.ProtocolID
	rmcResponse.MethodID = friends_wiiu.MethodRemoveFriend
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
