package notifications_3ds

import (
	"github.com/PretendoNetwork/friends/globals"
	nex "github.com/PretendoNetwork/nex-go"
	nintendo_notifications "github.com/PretendoNetwork/nex-protocols-go/nintendo-notifications"
	nintendo_notifications_types "github.com/PretendoNetwork/nex-protocols-go/nintendo-notifications/types"
)

func SendFriendshipCompleted(client *nex.PRUDPClient, friendPID uint32, senderPID *nex.PID) {
	notificationEvent := nintendo_notifications_types.NewNintendoNotificationEventGeneral()
	notificationEvent.U32Param = 0
	notificationEvent.U64Param1 = 0
	notificationEvent.U64Param2 = uint64(friendPID)

	eventObject := nintendo_notifications_types.NewNintendoNotificationEvent()
	eventObject.Type = 7
	eventObject.SenderPID = senderPID
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

	requestPacket, _ := nex.NewPRUDPPacketV0(client, nil)

	requestPacket.SetType(nex.DataPacket)
	requestPacket.AddFlag(nex.FlagNeedsAck)
	requestPacket.AddFlag(nex.FlagReliable)
	requestPacket.SetSourceStreamType(client.DestinationStreamType)
	requestPacket.SetSourcePort(client.DestinationPort)
	requestPacket.SetDestinationStreamType(client.SourceStreamType)
	requestPacket.SetDestinationPort(client.SourcePort)
	requestPacket.SetPayload(rmcRequestBytes)

	globals.SecureServer.Send(requestPacket)
}
