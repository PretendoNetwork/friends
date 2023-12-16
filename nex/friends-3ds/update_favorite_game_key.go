package nex_friends_3ds

import (
	database_3ds "github.com/PretendoNetwork/friends/database/3ds"
	"github.com/PretendoNetwork/friends/globals"
	notifications_3ds "github.com/PretendoNetwork/friends/notifications/3ds"
	nex "github.com/PretendoNetwork/nex-go"
	friends_3ds "github.com/PretendoNetwork/nex-protocols-go/friends-3ds"
	friends_3ds_types "github.com/PretendoNetwork/nex-protocols-go/friends-3ds/types"
)

func UpdateFavoriteGameKey(err error, packet nex.PacketInterface, callID uint32, gameKey *friends_3ds_types.GameKey) (*nex.RMCMessage, uint32) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.Errors.FPD.InvalidArgument
	}

	client := packet.Sender().(*nex.PRUDPClient)

	err = database_3ds.UpdateUserFavoriteGame(client.PID().LegacyValue(), gameKey)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return nil, nex.Errors.FPD.Unknown
	}

	go notifications_3ds.SendFavoriteUpdate(client, gameKey)

	rmcResponse := nex.NewRMCSuccess(globals.SecureServer, nil)
	rmcResponse.ProtocolID = friends_3ds.ProtocolID
	rmcResponse.MethodID = friends_3ds.MethodUpdateFavoriteGameKey
	rmcResponse.CallID = callID

	return rmcResponse, 0
}
