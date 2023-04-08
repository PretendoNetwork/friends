package nex_friends_3ds

import (
	database_3ds "github.com/PretendoNetwork/friends-secure/database/3ds"
	"github.com/PretendoNetwork/friends-secure/globals"
	notifications_3ds "github.com/PretendoNetwork/friends-secure/notifications/3ds"
	nex "github.com/PretendoNetwork/nex-go"
	friends_3ds "github.com/PretendoNetwork/nex-protocols-go/friends/3ds"
)

func UpdatePreference(err error, client *nex.Client, callID uint32, showOnline bool, showCurrentGame bool, showPlayedGame bool) {
	if !showCurrentGame {
		emptyPresence := friends_3ds.NewNintendoPresence()
		emptyPresence.GameKey = friends_3ds.NewGameKey()
		emptyPresence.ChangedFlags = 0xFFFFFFFF // All flags
		notifications_3ds.SendPresenceUpdate(client, emptyPresence)
	}
	if !showOnline {
		notifications_3ds.SendUserWentOfflineGlobally(client)
	}

	database_3ds.UpdateUserPreferences(client.PID(), showOnline, showCurrentGame)

	rmcResponse := nex.NewRMCResponse(friends_3ds.ProtocolID, callID)
	rmcResponse.SetSuccess(friends_3ds.MethodUpdatePreference, nil)

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
