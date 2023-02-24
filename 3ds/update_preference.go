package friends_3ds

import (
	database_3ds "github.com/PretendoNetwork/friends-secure/database/3ds"
	"github.com/PretendoNetwork/friends-secure/globals"
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func UpdatePreference(err error, client *nex.Client, callID uint32, showOnline bool, showCurrentGame bool, showPlayedGame bool) {
	if !showCurrentGame {
		emptyPresence := nexproto.NewNintendoPresence()
		emptyPresence.GameKey = nexproto.NewGameKey()
		emptyPresence.ChangedFlags = 4294967295 // FF FF FF FF, All flags
		sendPresenceUpdateNotification(client, emptyPresence)
	}
	if !showOnline {
		SendUserWentOfflineNotificationsGlobally(client)
	}

	database_3ds.UpdateUserPreferences(client.PID(), showOnline, showCurrentGame)

	rmcResponse := nex.NewRMCResponse(nexproto.Friends3DSProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.Friends3DSMethodUpdatePreference, nil)

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewPacketV0(client, nil)

	responsePacket.SetVersion(0)
	responsePacket.SetSource(0xA1)
	responsePacket.SetDestination(0xAF)
	responsePacket.SetType(nex.DataPacket)
	responsePacket.SetPayload(rmcResponseBytes)

	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)

	globals.NEXServer.Send(responsePacket)
}
