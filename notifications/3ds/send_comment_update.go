package notifications_3ds

import (
	"database/sql"

	database_3ds "github.com/PretendoNetwork/friends/database/3ds"
	"github.com/PretendoNetwork/friends/globals"
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/constants"
	"github.com/PretendoNetwork/nex-go/v2/types"
	friends_3ds_types "github.com/PretendoNetwork/nex-protocols-go/v2/friends-3ds/types"
	nintendo_notifications "github.com/PretendoNetwork/nex-protocols-go/v2/nintendo-notifications"
	nintendo_notifications_types "github.com/PretendoNetwork/nex-protocols-go/v2/nintendo-notifications/types"
)

func SendCommentUpdate(connection *nex.PRUDPConnection, comment string) {
	notificationEvent := nintendo_notifications_types.NewNintendoNotificationEventGeneral()
	notificationEvent.StrParam = types.NewString(comment)

	eventObject := nintendo_notifications_types.NewNintendoNotificationEvent()
	eventObject.Type = types.NewPrimitiveU32(3)
	eventObject.SenderPID = connection.PID()
	eventObject.DataHolder = types.NewAnyDataHolder()
	eventObject.DataHolder.TypeName = types.NewString("NintendoNotificationEventGeneral")
	eventObject.DataHolder.ObjectData = notificationEvent.Copy()

	stream := nex.NewByteStreamOut(globals.SecureEndpoint.LibraryVersions(), globals.SecureEndpoint.ByteStreamSettings())

	eventObject.WriteTo(stream)

	notificationRequest := nex.NewRMCRequest(globals.SecureEndpoint)
	notificationRequest.ProtocolID = nintendo_notifications.ProtocolID
	notificationRequest.CallID = 3810693103
	notificationRequest.MethodID = nintendo_notifications.MethodProcessNintendoNotificationEvent1
	notificationRequest.Parameters = stream.Bytes()

	notificationRequestBytes := notificationRequest.Bytes()

	friendsList, err := database_3ds.GetUserFriends(connection.PID().LegacyValue())
	if err != nil && err != sql.ErrNoRows {
		globals.Logger.Critical(err.Error())
	}

	if friendsList == nil {
		return
	}

	friendsList.Each(func(i int, friend *friends_3ds_types.FriendRelationship) bool {
		connectedUser := globals.ConnectedUsers[friend.PID.LegacyValue()]

		if connectedUser != nil {
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

		return false
	})
}
