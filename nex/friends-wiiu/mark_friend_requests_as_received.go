package nex_friends_wiiu

import (
	database_wiiu "github.com/PretendoNetwork/friends/database/wiiu"
	"github.com/PretendoNetwork/friends/globals"
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	friends_wiiu "github.com/PretendoNetwork/nex-protocols-go/v2/friends-wiiu"
)

func MarkFriendRequestsAsReceived(err error, packet nex.PacketInterface, callID uint32, ids *types.List[*types.PrimitiveU64]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.FPD.InvalidArgument, "") // TODO - Add error message
	}

	if ids.Each(func(i int, id *types.PrimitiveU64) bool {
		err = database_wiiu.SetFriendRequestReceived(id.Value)
		if err != nil {
			globals.Logger.Critical(err.Error())
			return true
		}

		return false
	}) {
		return nil, nex.NewError(nex.ResultCodes.FPD.Unknown, "") // TODO - Add error message
	}

	rmcResponse := nex.NewRMCSuccess(globals.SecureEndpoint, nil)
	rmcResponse.ProtocolID = friends_wiiu.ProtocolID
	rmcResponse.MethodID = friends_wiiu.MethodMarkFriendRequestsAsReceived
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
