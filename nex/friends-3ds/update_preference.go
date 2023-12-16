package nex_friends_3ds

import (
	database_3ds "github.com/PretendoNetwork/friends/database/3ds"
	"github.com/PretendoNetwork/friends/globals"
	notifications_3ds "github.com/PretendoNetwork/friends/notifications/3ds"
	nex "github.com/PretendoNetwork/nex-go"
	friends_3ds "github.com/PretendoNetwork/nex-protocols-go/friends-3ds"
	friends_3ds_types "github.com/PretendoNetwork/nex-protocols-go/friends-3ds/types"
)

func UpdatePreference(err error, packet nex.PacketInterface, callID uint32, showOnline bool, showCurrentGame bool, showPlayedGame bool) (*nex.RMCMessage, uint32) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.Errors.FPD.InvalidArgument
	}

	client := packet.Sender().(*nex.PRUDPClient)

	err = database_3ds.UpdateUserPreferences(client.PID().LegacyValue(), showOnline, showCurrentGame)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return nil, nex.Errors.FPD.Unknown
	}

	if !showCurrentGame {
		emptyPresence := friends_3ds_types.NewNintendoPresence()
		emptyPresence.GameKey = friends_3ds_types.NewGameKey()
		emptyPresence.ChangedFlags = 0xFFFFFFFF // All flags
		notifications_3ds.SendPresenceUpdate(client, emptyPresence)
	}
	if !showOnline {
		notifications_3ds.SendUserWentOfflineGlobally(client)
	}

	rmcResponse := nex.NewRMCSuccess(globals.SecureServer, nil)
	rmcResponse.ProtocolID = friends_3ds.ProtocolID
	rmcResponse.MethodID = friends_3ds.MethodUpdatePreference
	rmcResponse.CallID = callID

	return rmcResponse, 0
}
