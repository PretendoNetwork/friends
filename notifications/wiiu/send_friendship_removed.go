package notifications_wiiu

import (
	"github.com/PretendoNetwork/friends/globals"
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/constants"
	"github.com/PretendoNetwork/nex-go/v2/types"
	nintendo_notifications "github.com/PretendoNetwork/nex-protocols-go/v2/nintendo-notifications"
	nintendo_notifications_types "github.com/PretendoNetwork/nex-protocols-go/v2/nintendo-notifications/types"
)

func SendFriendshipRemoved(connection *nex.PRUDPConnection, senderPID *types.PID) {
	nintendoNotificationEventGeneral := nintendo_notifications_types.NewNintendoNotificationEventGeneral()

	eventObject := nintendo_notifications_types.NewNintendoNotificationEvent()
	eventObject.Type = types.NewPrimitiveU32(26)
	eventObject.SenderPID = senderPID.Copy().(*types.PID)
	eventObject.DataHolder = types.NewAnyDataHolder()
	eventObject.DataHolder.TypeName = types.NewString("NintendoNotificationEventGeneral")
	eventObject.DataHolder.ObjectData = nintendoNotificationEventGeneral.Copy()

	stream := nex.NewByteStreamOut(globals.SecureEndpoint.LibraryVersions(), globals.SecureEndpoint.ByteStreamSettings())

	eventObject.WriteTo(stream)

	notificationRequest := nex.NewRMCRequest(globals.SecureEndpoint)
	notificationRequest.ProtocolID = nintendo_notifications.ProtocolID
	notificationRequest.CallID = 3810693103
	notificationRequest.MethodID = nintendo_notifications.MethodProcessNintendoNotificationEvent1
	notificationRequest.Parameters = stream.Bytes()

	notificationRequestBytes := notificationRequest.Bytes()

	requestPacket, _ := nex.NewPRUDPPacketV0(globals.SecureEndpoint.Server, connection, nil)

	requestPacket.SetType(constants.DataPacket)
	requestPacket.AddFlag(constants.PacketFlagNeedsAck)
	requestPacket.AddFlag(constants.PacketFlagReliable)
	requestPacket.SetSourceVirtualPortStreamType(connection.StreamType)
	requestPacket.SetSourceVirtualPortStreamID(globals.SecureEndpoint.StreamID)
	requestPacket.SetDestinationVirtualPortStreamType(connection.StreamType)
	requestPacket.SetDestinationVirtualPortStreamID(connection.StreamID)
	requestPacket.SetPayload(notificationRequestBytes)

	globals.SecureServer.Send(requestPacket)
}
