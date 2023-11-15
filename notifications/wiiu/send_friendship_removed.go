package notifications_wiiu

import (
	"github.com/PretendoNetwork/friends/globals"
	nex "github.com/PretendoNetwork/nex-go"
	nintendo_notifications "github.com/PretendoNetwork/nex-protocols-go/nintendo-notifications"
	nintendo_notifications_types "github.com/PretendoNetwork/nex-protocols-go/nintendo-notifications/types"
)

func SendFriendshipRemoved(client *nex.PRUDPClient, senderPID *nex.PID) {
	nintendoNotificationEventGeneral := nintendo_notifications_types.NewNintendoNotificationEventGeneral()

	eventObject := nintendo_notifications_types.NewNintendoNotificationEvent()
	eventObject.Type = 26
	eventObject.SenderPID = senderPID
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
