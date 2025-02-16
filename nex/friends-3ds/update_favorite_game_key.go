package nex_friends_3ds

import (
	database_3ds "github.com/PretendoNetwork/friends/database/3ds"
	"github.com/PretendoNetwork/friends/globals"
	notifications_3ds "github.com/PretendoNetwork/friends/notifications/3ds"
	nex "github.com/PretendoNetwork/nex-go/v2"
	friends_3ds "github.com/PretendoNetwork/nex-protocols-go/v2/friends-3ds"
	friends_3ds_types "github.com/PretendoNetwork/nex-protocols-go/v2/friends-3ds/types"
)

func UpdateFavoriteGameKey(err error, packet nex.PacketInterface, callID uint32, gameKey friends_3ds_types.GameKey) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.FPD.InvalidArgument, "") // TODO - Add error message
	}

	connection := packet.Sender().(*nex.PRUDPConnection)

	err = database_3ds.UpdateUserFavoriteGame(uint32(connection.PID()), gameKey)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return nil, nex.NewError(nex.ResultCodes.FPD.Unknown, "") // TODO - Add error message
	}

	go notifications_3ds.SendFavoriteUpdate(connection, gameKey)

	rmcResponse := nex.NewRMCSuccess(globals.SecureEndpoint, nil)
	rmcResponse.ProtocolID = friends_3ds.ProtocolID
	rmcResponse.MethodID = friends_3ds.MethodUpdateFavoriteGameKey
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
