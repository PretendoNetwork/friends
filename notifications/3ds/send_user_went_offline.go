package notifications_3ds

import (
	"database/sql"

	database_3ds "github.com/CloudnetworkTeam/friends/database/3ds"
	"github.com/CloudnetworkTeam/friends/globals"
	nex "github.com/PretendoNetwork/nex-go"
	nintendo_notifications "github.com/PretendoNetwork/nex-protocols-go/nintendo-notifications"
	nintendo_notifications_types "github.com/PretendoNetwork/nex-protocols-go/nintendo-notifications/types"
)

func SendUserWentOfflineGlobally(client *nex.Client) {
	friendsList, err := database_3ds.GetUserFriends(client.PID())
	if err != nil && err != sql.ErrNoRows {
		globals.Logger.Critical(err.Error())
	}

	for i := 0; i < len(friendsList); i++ {
		SendUserWentOffline(client, friendsList[i].PID)
	}
}

func SendUserWentOffline(client *nex.Client, pid uint32) {
	notificationEvent := nintendo_notifications_types.NewNintendoNotificationEventGeneral()

	eventObject := nintendo_notifications_types.NewNintendoNotificationEvent()
	eventObject.Type = 10
	eventObject.SenderPID = client.PID()
	eventObject.DataHolder = nex.NewDataHolder()
	eventObject.DataHolder.SetTypeName("NintendoNotificationEventGeneral")
	eventObject.DataHolder.SetObjectData(notificationEvent)

	stream := nex.NewStreamOut(globals.SecureServer)
	eventObjectBytes := eventObject.Bytes(stream)

	rmcRequest := nex.NewRMCRequest()
	rmcRequest.SetProtocolID(nintendo_notifications.ProtocolID)
	rmcRequest.SetCallID(3810693103)
	rmcRequest.SetMethodID(nintendo_notifications.MethodProcessNintendoNotificationEvent1)
	rmcRequest.SetParameters(eventObjectBytes)

	rmcRequestBytes := rmcRequest.Bytes()

	connectedUser := globals.ConnectedUsers[pid]

	if connectedUser != nil {
		requestPacket, _ := nex.NewPacketV0(connectedUser.Client, nil)

		requestPacket.SetVersion(0)
		requestPacket.SetSource(0xA1)
		requestPacket.SetDestination(0xAF)
		requestPacket.SetType(nex.DataPacket)
		requestPacket.SetPayload(rmcRequestBytes)

		requestPacket.AddFlag(nex.FlagNeedsAck)
		requestPacket.AddFlag(nex.FlagReliable)

		globals.SecureServer.Send(requestPacket)
	}
}
