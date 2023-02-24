package notifications_3ds

import (
	database_3ds "github.com/PretendoNetwork/friends-secure/database/3ds"
	"github.com/PretendoNetwork/friends-secure/globals"
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func SendCommentUpdate(client *nex.Client, comment string) {
	notificationEvent := nexproto.NewNintendoNotificationEventGeneral()
	notificationEvent.StrParam = comment

	eventObject := nexproto.NewNintendoNotificationEvent()
	eventObject.Type = 3
	eventObject.SenderPID = client.PID()
	eventObject.DataHolder = nex.NewDataHolder()
	eventObject.DataHolder.SetTypeName("NintendoNotificationEventGeneral")
	eventObject.DataHolder.SetObjectData(notificationEvent)

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
