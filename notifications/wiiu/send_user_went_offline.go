package notifications_wiiu

import (
	"time"

	"github.com/PretendoNetwork/friends/database"
	database_wiiu "github.com/PretendoNetwork/friends/database/wiiu"
	"github.com/PretendoNetwork/friends/globals"
	nex "github.com/PretendoNetwork/nex-go"
	nintendo_notifications "github.com/PretendoNetwork/nex-protocols-go/nintendo-notifications"
	nintendo_notifications_types "github.com/PretendoNetwork/nex-protocols-go/nintendo-notifications/types"
)

func SendUserWentOfflineGlobally(client *nex.PRUDPClient) {
	friendsList, err := database_wiiu.GetUserFriendList(client.PID())
	if err != nil && err != database.ErrEmptyList {
		globals.Logger.Critical(err.Error())
	}

	for i := 0; i < len(friendsList); i++ {
		SendUserWentOffline(client, friendsList[i].NNAInfo.PrincipalBasicInfo.PID)
	}
}

func SendUserWentOffline(client *nex.PRUDPClient, pid uint32) {
	lastOnline := nex.NewDateTime(0)
	lastOnline.FromTimestamp(time.Now())

	nintendoNotificationEventGeneral := nintendo_notifications_types.NewNintendoNotificationEventGeneral()

	nintendoNotificationEventGeneral.U32Param = 0
	nintendoNotificationEventGeneral.U64Param1 = 0
	nintendoNotificationEventGeneral.U64Param2 = lastOnline.Value()
	nintendoNotificationEventGeneral.StrParam = ""

	eventObject := nintendo_notifications_types.NewNintendoNotificationEvent()
	eventObject.Type = 10
	eventObject.SenderPID = client.PID()
	eventObject.DataHolder = nex.NewDataHolder()
	eventObject.DataHolder.SetTypeName("NintendoNotificationEventGeneral")
	eventObject.DataHolder.SetObjectData(nintendoNotificationEventGeneral)

	stream := nex.NewStreamOut(globals.SecureServer)
	stream.WriteStructure(eventObject)

	rmcRequest := nex.NewRMCRequest()
	rmcRequest.ProtocolID = nintendo_notifications.ProtocolID
	rmcRequest.CallID = 3810693103
	rmcRequest.MethodID = nintendo_notifications.MethodProcessNintendoNotificationEvent1
	rmcRequest.Parameters = stream.Bytes()

	rmcRequestBytes := rmcRequest.Bytes()

	connectedUser := globals.ConnectedUsers[pid]

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
