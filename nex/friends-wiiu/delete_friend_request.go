package nex_friends_wiiu

import (
	"github.com/PretendoNetwork/friends/database"
	database_wiiu "github.com/PretendoNetwork/friends/database/wiiu"
	"github.com/PretendoNetwork/friends/globals"
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	friends_wiiu "github.com/PretendoNetwork/nex-protocols-go/v2/friends-wiiu"
)

func DeleteFriendRequest(err error, packet nex.PacketInterface, callID uint32, id types.UInt64) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.FPD.InvalidArgument, "") // TODO - Add error message
	}

	err = database_wiiu.SetFriendRequestDenied(uint64(id))
	if err != nil {
		if err == database.ErrFriendRequestNotFound {
			return nil, nex.NewError(nex.ResultCodes.FPD.InvalidMessageID, "") // TODO - Add error message
		} else {
			globals.Logger.Critical(err.Error())
			return nil, nex.NewError(nex.ResultCodes.FPD.Unknown, "") // TODO - Add error message
		}
	}

	rmcResponse := nex.NewRMCSuccess(globals.SecureEndpoint, nil)
	rmcResponse.ProtocolID = friends_wiiu.ProtocolID
	rmcResponse.MethodID = friends_wiiu.MethodDeleteFriendRequest
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
