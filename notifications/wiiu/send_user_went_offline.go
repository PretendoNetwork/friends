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

func SendUserWentOfflineGlobally(client *nex.Client) {
	friendsList, err := database_wiiu.GetUserFriendList(client.PID())
	if err != nil && err != database.ErrEmptyList {
		globals.Logger.Critical(err.Error())
	}

	for i := 0; i < len(friendsList); i++ {
		SendUserWentOffline(client, friendsList[i].NNAInfo.PrincipalBasicInfo.PID)
	}
}

func SendUserWentOffline(client *nex.Client, pid uint32) {
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
	rmcRequest.SetProtocolID(nintendo_notifications.ProtocolID)
	rmcRequest.SetCallID(3810693103)
	rmcRequest.SetMethodID(nintendo_notifications.MethodProcessNintendoNotificationEvent1)
	rmcRequest.SetParameters(stream.Bytes())

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
