package nex_friends_3ds

import (
	database_3ds "github.com/PretendoNetwork/friends/database/3ds"
	"github.com/PretendoNetwork/friends/globals"
	notifications_3ds "github.com/PretendoNetwork/friends/notifications/3ds"
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	friends_3ds "github.com/PretendoNetwork/nex-protocols-go/v2/friends-3ds"
	friends_3ds_types "github.com/PretendoNetwork/nex-protocols-go/v2/friends-3ds/types"
)

func UpdatePreference(err error, packet nex.PacketInterface, callID uint32, publicMode *types.PrimitiveBool, showGame *types.PrimitiveBool, showPlayedGame *types.PrimitiveBool) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.FPD.InvalidArgument, "") // TODO - Add error message
	}

	connection := packet.Sender().(*nex.PRUDPConnection)

	err = database_3ds.UpdateUserPreferences(connection.PID().LegacyValue(), publicMode.Value, showGame.Value)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return nil, nex.NewError(nex.ResultCodes.FPD.Unknown, "") // TODO - Add error message
	}

	if !showGame.Value {
		emptyPresence := friends_3ds_types.NewNintendoPresence()
		emptyPresence.GameKey = friends_3ds_types.NewGameKey()
		emptyPresence.ChangedFlags = types.NewPrimitiveU32(0xFFFFFFFF) // * All flags
		notifications_3ds.SendPresenceUpdate(connection, emptyPresence)
	}

	if !publicMode.Value {
		notifications_3ds.SendUserWentOfflineGlobally(connection)
	}

	rmcResponse := nex.NewRMCSuccess(globals.SecureEndpoint, nil)
	rmcResponse.ProtocolID = friends_3ds.ProtocolID
	rmcResponse.MethodID = friends_3ds.MethodUpdatePreference
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
