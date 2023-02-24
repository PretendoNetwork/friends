package friends_3ds

import (
	database_3ds "github.com/PretendoNetwork/friends-secure/database/3ds"
	"github.com/PretendoNetwork/friends-secure/globals"
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func UpdatePresence(err error, client *nex.Client, callID uint32, presence *nexproto.NintendoPresence, showGame bool) {
	currentPresence := presence

	// Send an entirely empty status, with every flag set to update
	if !showGame {
		currentPresence = nexproto.NewNintendoPresence()
		currentPresence.GameKey = nexproto.NewGameKey()
		currentPresence.ChangedFlags = 4294967295 // FF FF FF FF, All flags
	}

	go sendPresenceUpdateNotification(client, currentPresence)
	globals.ConnectedUsers[client.PID()].Presence = currentPresence

	rmcResponse := nex.NewRMCResponse(nexproto.Friends3DSProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.Friends3DSMethodUpdatePresence, nil)

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

func sendPresenceUpdateNotification(client *nex.Client, presence *nexproto.NintendoPresence) {
	eventObject := nexproto.NewNintendoNotificationEvent()
	eventObject.Type = 1
	eventObject.SenderPID = client.PID()
	eventObject.DataHolder = nex.NewDataHolder()
	eventObject.DataHolder.SetTypeName("NintendoPresence")
	eventObject.DataHolder.SetObjectData(presence)

	stream := nex.NewStreamOut(globals.NEXServer)
	eventObjectBytes := eventObject.Bytes(stream)

	rmcRequest := nex.NewRMCRequest()
	rmcRequest.SetProtocolID(nexproto.NintendoNotificationsProtocolID)
	rmcRequest.SetCallID(3810693103)
	rmcRequest.SetMethodID(nexproto.NintendoNotificationsMethodProcessNintendoNotificationEvent1)
	rmcRequest.SetParameters(eventObjectBytes)

	rmcRequestBytes := rmcRequest.Bytes()

	friendsList := database_3ds.GetUserFriends(client.PID())

	for i := 0; i < len(friendsList); i++ {

		connectedUser := globals.ConnectedUsers[friendsList[i].PID]

		if connectedUser != nil {

			requestPacket, _ := nex.NewPacketV0(connectedUser.Client, nil)

			requestPacket.SetVersion(0)
			requestPacket.SetSource(0xA1)
			requestPacket.SetDestination(0xAF)
			requestPacket.SetType(nex.DataPacket)
			requestPacket.SetPayload(rmcRequestBytes)

			requestPacket.AddFlag(nex.FlagNeedsAck)
			requestPacket.AddFlag(nex.FlagReliable)

			globals.NEXServer.Send(requestPacket)
		}
	}
}
