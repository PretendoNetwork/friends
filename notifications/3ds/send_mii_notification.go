package notifications_3ds

import (
	"database/sql"

	database_3ds "github.com/PretendoNetwork/friends/database/3ds"
	"github.com/PretendoNetwork/friends/globals"
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/constants"
	"github.com/PretendoNetwork/nex-go/v2/types"
	nintendo_notifications "github.com/PretendoNetwork/nex-protocols-go/v2/nintendo-notifications"
	nintendo_notifications_types "github.com/PretendoNetwork/nex-protocols-go/v2/nintendo-notifications/types"
)

func SendMiiUpdateNotification(connection *nex.PRUDPConnection) {
	notificationEvent := nintendo_notifications_types.NewNintendoNotificationEventGeneral()

	eventObject := nintendo_notifications_types.NewNintendoNotificationEvent()
	eventObject.Type = types.NewUInt32(5)
	eventObject.SenderPID = connection.PID()
	eventObject.DataHolder = types.NewDataHolder()
	eventObject.DataHolder.Object = notificationEvent.Copy().(nintendo_notifications_types.NintendoNotificationEventGeneral)

	stream := nex.NewByteStreamOut(globals.SecureEndpoint.LibraryVersions(), globals.SecureEndpoint.ByteStreamSettings())

	eventObject.WriteTo(stream)

	notificationRequest := nex.NewRMCRequest(globals.SecureEndpoint)
	notificationRequest.ProtocolID = nintendo_notifications.ProtocolID
	notificationRequest.CallID = 3810693103
	notificationRequest.MethodID = nintendo_notifications.MethodProcessNintendoNotificationEvent1
	notificationRequest.Parameters = stream.Bytes()

	notificationRequestBytes := notificationRequest.Bytes()

	friendsList, err := database_3ds.GetUserFriends(uint32(connection.PID()))
	if err != nil && err != sql.ErrNoRows {
		globals.Logger.Critical(err.Error())
	}

	if friendsList == nil {
		return
	}

	for _, friend := range friendsList {
		connectedUser, ok := globals.ConnectedUsers.Get(uint32(friend.PID))

		if ok && connectedUser != nil {
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
}
