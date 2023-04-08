package notifications_wiiu

import (
	"fmt"

	database_wiiu "github.com/PretendoNetwork/friends-secure/database/wiiu"
	"github.com/PretendoNetwork/friends-secure/globals"
	nex "github.com/PretendoNetwork/nex-go"
	friends_wiiu "github.com/PretendoNetwork/nex-protocols-go/friends/wiiu"
	nintendo_notifications "github.com/PretendoNetwork/nex-protocols-go/nintendo-notifications"
)

func SendPresenceUpdate(presence *friends_wiiu.NintendoPresenceV2) {
	eventObject := nintendo_notifications.NewNintendoNotificationEvent()
	eventObject.Type = 24
	eventObject.SenderPID = presence.PID
	eventObject.DataHolder = nex.NewDataHolder()
	eventObject.DataHolder.SetTypeName("NintendoPresenceV2")
	eventObject.DataHolder.SetObjectData(presence)

	stream := nex.NewStreamOut(globals.NEXServer)
	eventObjectBytes := eventObject.Bytes(stream)

	rmcRequest := nex.NewRMCRequest()
	rmcRequest.SetProtocolID(nintendo_notifications.ProtocolID)
	rmcRequest.SetCallID(3810693103)
	rmcRequest.SetMethodID(nintendo_notifications.MethodProcessNintendoNotificationEvent2)
	rmcRequest.SetParameters(eventObjectBytes)

	rmcRequestBytes := rmcRequest.Bytes()

	friendList := database_wiiu.GetUserFriendList(presence.PID)

	for i := 0; i < len(friendList); i++ {
		if friendList[i] == nil || friendList[i].NNAInfo == nil || friendList[i].NNAInfo.PrincipalBasicInfo == nil {
			// TODO: Fix this
			pid := presence.PID
			var friendPID uint32 = 0

			if friendList[i] != nil && friendList[i].Presence != nil {
				// TODO: Better track the bad users PID
				friendPID = friendList[i].Presence.PID
			}

			globals.Logger.Error(fmt.Sprintf("User %d has friend %d with bad presence data", pid, friendPID))

			if friendList[i] == nil {
				globals.Logger.Error(fmt.Sprintf("%d friendList[i] nil", friendPID))
			} else if friendList[i].NNAInfo == nil {
				globals.Logger.Error(fmt.Sprintf("%d friendList[i].NNAInfo is nil", friendPID))
			} else if friendList[i].NNAInfo.PrincipalBasicInfo == nil {
				globals.Logger.Error(fmt.Sprintf("%d friendList[i].NNAInfo.PrincipalBasicInfo is nil", friendPID))
			}

			continue
		}

		friendPID := friendList[i].NNAInfo.PrincipalBasicInfo.PID
		connectedUser := globals.ConnectedUsers[friendPID]

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
