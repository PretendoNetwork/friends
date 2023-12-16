package notifications_3ds

import (
	"database/sql"

	database_3ds "github.com/PretendoNetwork/friends/database/3ds"
	"github.com/PretendoNetwork/friends/globals"
	nex "github.com/PretendoNetwork/nex-go"
	nintendo_notifications "github.com/PretendoNetwork/nex-protocols-go/nintendo-notifications"
	nintendo_notifications_types "github.com/PretendoNetwork/nex-protocols-go/nintendo-notifications/types"
)

func SendMiiUpdateNotification(client *nex.PRUDPClient) {
	notificationEvent := nintendo_notifications_types.NewNintendoNotificationEventGeneral()

	eventObject := nintendo_notifications_types.NewNintendoNotificationEvent()
	eventObject.Type = 5
	eventObject.SenderPID = client.PID()
	eventObject.DataHolder = nex.NewDataHolder()
	eventObject.DataHolder.SetTypeName("NintendoNotificationEventGeneral")
	eventObject.DataHolder.SetObjectData(notificationEvent)

	stream := nex.NewStreamOut(globals.SecureServer)
	eventObjectBytes := eventObject.Bytes(stream)

	rmcRequest := nex.NewRMCRequest(globals.SecureServer)
	rmcRequest.ProtocolID = nintendo_notifications.ProtocolID
	rmcRequest.CallID = 3810693103
	rmcRequest.MethodID = nintendo_notifications.MethodProcessNintendoNotificationEvent1
	rmcRequest.Parameters = eventObjectBytes

	rmcRequestBytes := rmcRequest.Bytes()

	friendsList, err := database_3ds.GetUserFriends(client.PID().LegacyValue())
	if err != nil && err != sql.ErrNoRows {
		globals.Logger.Critical(err.Error())
	}

	for i := 0; i < len(friendsList); i++ {

		connectedUser := globals.ConnectedUsers[friendsList[i].PID.LegacyValue()]

		if connectedUser != nil {
			requestPacket, _ := nex.NewPRUDPPacketV0(connectedUser.Client, nil)

			requestPacket.SetType(nex.DataPacket)
			requestPacket.AddFlag(nex.FlagNeedsAck)
			requestPacket.AddFlag(nex.FlagReliable)
			requestPacket.SetSourceStreamType(connectedUser.Client.DestinationStreamType)
			requestPacket.SetSourcePort(connectedUser.Client.DestinationPort)
			requestPacket.SetDestinationStreamType(connectedUser.Client.SourceStreamType)
			requestPacket.SetDestinationPort(connectedUser.Client.SourcePort)
			requestPacket.SetPayload(rmcRequestBytes)

			globals.SecureServer.Send(requestPacket)
		}
	}
}
