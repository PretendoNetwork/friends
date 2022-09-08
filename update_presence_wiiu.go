package main

import (
	"fmt"

	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func updatePresenceWiiU(err error, client *nex.Client, callID uint32, presence *nexproto.NintendoPresenceV2) {
	pid := client.PID()

	presence.Online = true      // Force online status. I have no idea why this is always false
	presence.PID = client.PID() // WHY IS THIS SET TO 0 BY DEFAULT??

	connectedUsers[pid].Presence = presence
	sendUpdatePresenceWiiUNotifications(presence)

	rmcResponse := nex.NewRMCResponse(nexproto.FriendsWiiUProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.FriendsWiiUMethodUpdatePresence, nil)

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
	eventObject.DataHolder.SetTypeName("NintendoPresenceV2")
	eventObject.DataHolder.SetObjectData(presence)

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
		if friendList[i] == nil || friendList[i].NNAInfo == nil || friendList[i].NNAInfo.PrincipalBasicInfo == nil {
			// TODO: Fix this
			pid := presence.PID
			var friendPID uint32 = 0

			if friendList[i] != nil && friendList[i].Presence != nil {
				// TODO: Better track the bad users PID
				friendPID = friendList[i].Presence.PID
			}

			logger.Error(fmt.Sprintf("User %d has friend %d with bad presence data", pid, friendPID))

			if friendList[i] == nil {
				logger.Error(fmt.Sprintf("%d friendList[i] nil", friendPID))
			} else if friendList[i].NNAInfo == nil {
				logger.Error(fmt.Sprintf("%d friendList[i].NNAInfo is nil", friendPID))
			} else if friendList[i].NNAInfo.PrincipalBasicInfo == nil {
				logger.Error(fmt.Sprintf("%d friendList[i].NNAInfo.PrincipalBasicInfo is nil", friendPID))
			}

			continue
		}

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
