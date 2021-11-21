package main

import (
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func updatePresenceWiiU(err error, client *nex.Client, callID uint32, presence *nexproto.NintendoPresenceV2) {
	pid := client.PID()

	presence.Online = true      // Force online status. I have no idea why this is always false
	presence.PID = client.PID() // WHY IS THIS SET TO 0 BY DEFAULT??

	connectedUsers[pid].Presence = presence
	sendUpdatePresenceWiiUNotifications(presence)

	rmcResponse := nex.NewRMCResponse(nexproto.FriendsProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.FriendsMethodUpdatePresence, nil)

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewPacketV0(client, nil)

	responsePacket.SetVersion(0)
	responsePacket.SetSource(0xA1)
	responsePacket.SetDestination(0xAF)
	responsePacket.SetType(nex.DataPacket)
	responsePacket.SetPayload(rmcResponseBytes)

	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)

	nexServer.Send(responsePacket)
}

func sendUpdatePresenceWiiUNotifications(presence *nexproto.NintendoPresenceV2) {
	eventObject := nexproto.NewNintendoNotificationEvent()
	eventObject.Type = 24
	eventObject.SenderPID = presence.PID
	eventObject.DataHolder = nex.NewDataHolder()
	eventObject.DataHolder.Name = "NintendoPresenceV2"
	eventObject.DataHolder.Object = presence

	stream := nex.NewStreamOut(nexServer)
	eventObjectBytes := eventObject.Bytes(stream)

	rmcRequest := nex.NewRMCRequest()
	rmcRequest.SetProtocolID(nexproto.NintendoNotificationsProtocolID)
	rmcRequest.SetCallID(3810693103)
	rmcRequest.SetMethodID(nexproto.NintendoNotificationsMethodProcessNintendoNotificationEvent2)
	rmcRequest.SetParameters(eventObjectBytes)

	rmcRequestBytes := rmcRequest.Bytes()

	friendList := getUserFriendList(presence.PID)

	for i := 0; i < len(friendList); i++ {
		friendPID := friendList[i].NNAInfo.PrincipalBasicInfo.PID
		connectedUser := connectedUsers[friendPID]

		if connectedUser != nil {
			requestPacket, _ := nex.NewPacketV0(connectedUser.Client, nil)

			requestPacket.SetVersion(0)
			requestPacket.SetSource(0xA1)
			requestPacket.SetDestination(0xAF)
			requestPacket.SetType(nex.DataPacket)
			requestPacket.SetPayload(rmcRequestBytes)

			requestPacket.AddFlag(nex.FlagNeedsAck)
			requestPacket.AddFlag(nex.FlagReliable)

			nexServer.Send(requestPacket)
		}
	}
}
