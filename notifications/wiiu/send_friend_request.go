package notifications_wiiu

import (
	"github.com/PretendoNetwork/friends/globals"
	nex "github.com/PretendoNetwork/nex-go"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/friends-wiiu/types"
	nintendo_notifications "github.com/PretendoNetwork/nex-protocols-go/nintendo-notifications"
	nintendo_notifications_types "github.com/PretendoNetwork/nex-protocols-go/nintendo-notifications/types"
)

func SendFriendRequest(client *nex.PRUDPClient, friendRequestNotificationData *friends_wiiu_types.FriendRequest) {
	eventObject := nintendo_notifications_types.NewNintendoNotificationEvent()
	eventObject.Type = 27
	eventObject.SenderPID = friendRequestNotificationData.PrincipalInfo.PID
	eventObject.DataHolder = nex.NewDataHolder()
	eventObject.DataHolder.SetTypeName("FriendRequest")
	eventObject.DataHolder.SetObjectData(friendRequestNotificationData)

	stream := nex.NewStreamOut(globals.SecureServer)
	eventObjectBytes := eventObject.Bytes(stream)

	rmcRequest := nex.NewRMCRequest(globals.SecureServer)
	rmcRequest.ProtocolID = nintendo_notifications.ProtocolID
	rmcRequest.CallID = 3810693103
	rmcRequest.MethodID = nintendo_notifications.MethodProcessNintendoNotificationEvent2
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
