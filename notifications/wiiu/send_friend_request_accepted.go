package notifications_wiiu

import (
	"github.com/PretendoNetwork/friends/globals"
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/constants"
	"github.com/PretendoNetwork/nex-go/v2/types"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/v2/friends-wiiu/types"
	nintendo_notifications "github.com/PretendoNetwork/nex-protocols-go/v2/nintendo-notifications"
	nintendo_notifications_types "github.com/PretendoNetwork/nex-protocols-go/v2/nintendo-notifications/types"
)

func SendFriendRequestAccepted(connection *nex.PRUDPConnection, friendInfo friends_wiiu_types.FriendInfo) {
	eventObject := nintendo_notifications_types.NewNintendoNotificationEvent()
	eventObject.Type = types.NewUInt32(30)
	eventObject.SenderPID = friendInfo.NNAInfo.PrincipalBasicInfo.PID.Copy().(types.PID)
	eventObject.DataHolder = types.NewDataHolder()
	eventObject.DataHolder.Object = friendInfo.Copy().(friends_wiiu_types.FriendInfo)

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
