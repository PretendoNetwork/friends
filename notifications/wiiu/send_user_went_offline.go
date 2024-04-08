package notifications_wiiu

import (
	"github.com/PretendoNetwork/friends/database"
	database_wiiu "github.com/PretendoNetwork/friends/database/wiiu"
	"github.com/PretendoNetwork/friends/globals"
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/constants"
	"github.com/PretendoNetwork/nex-go/v2/types"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/v2/friends-wiiu/types"
	nintendo_notifications "github.com/PretendoNetwork/nex-protocols-go/v2/nintendo-notifications"
	nintendo_notifications_types "github.com/PretendoNetwork/nex-protocols-go/v2/nintendo-notifications/types"
)

func SendUserWentOfflineGlobally(connection *nex.PRUDPConnection) {
	friendsList, err := database_wiiu.GetUserFriendList(connection.PID().LegacyValue())
	if err != nil && err != database.ErrEmptyList {
		globals.Logger.Critical(err.Error())
	}

	if friendsList == nil {
		return
	}

	friendsList.Each(func(i int, friend *friends_wiiu_types.FriendInfo) bool {
		SendUserWentOffline(connection, friend.NNAInfo.PrincipalBasicInfo.PID)
		return false
	})
}

func SendUserWentOffline(connection *nex.PRUDPConnection, pid *types.PID) {
	lastOnline := types.NewDateTime(0).Now()

	nintendoNotificationEventGeneral := nintendo_notifications_types.NewNintendoNotificationEventGeneral()

	nintendoNotificationEventGeneral.U32Param = types.NewPrimitiveU32(0)
	nintendoNotificationEventGeneral.U64Param1 = types.NewPrimitiveU64(0)
	nintendoNotificationEventGeneral.U64Param2 = types.NewPrimitiveU64(lastOnline.Value())
	nintendoNotificationEventGeneral.StrParam = types.NewString("")

	eventObject := nintendo_notifications_types.NewNintendoNotificationEvent()
	eventObject.Type = types.NewPrimitiveU32(10)
	eventObject.SenderPID = connection.PID().Copy().(*types.PID)
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

	connectedUser := globals.ConnectedUsers[pid.LegacyValue()]

	if connectedUser != nil {
		requestPacket, _ := nex.NewPRUDPPacketV0(globals.SecureEndpoint.Server, connectedUser.Connection, nil)

		requestPacket.SetType(constants.DataPacket)
		requestPacket.AddFlag(constants.PacketFlagNeedsAck)
		requestPacket.AddFlag(constants.PacketFlagReliable)
		requestPacket.SetSourceVirtualPortStreamType(connectedUser.Connection.StreamType)
		requestPacket.SetSourceVirtualPortStreamID(globals.SecureEndpoint.StreamID)
		requestPacket.SetDestinationVirtualPortStreamType(connectedUser.Connection.StreamType)
		requestPacket.SetDestinationVirtualPortStreamID(connectedUser.Connection.StreamID)
		requestPacket.SetPayload(notificationRequestBytes)

		globals.SecureServer.Send(requestPacket)
	}
}
