package nex_friends_wiiu

import (
	"time"

	"github.com/PretendoNetwork/friends/database"
	database_wiiu "github.com/PretendoNetwork/friends/database/wiiu"
	"github.com/PretendoNetwork/friends/globals"
	notifications_wiiu "github.com/PretendoNetwork/friends/notifications/wiiu"
	"github.com/PretendoNetwork/friends/utility"
	nex "github.com/PretendoNetwork/nex-go"
	friends_wiiu "github.com/PretendoNetwork/nex-protocols-go/friends-wiiu"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/friends-wiiu/types"
)

func AddFriendRequest(err error, packet nex.PacketInterface, callID uint32, pid uint32, unknown2 uint8, message string, unknown4 uint8, unknown5 string, gameKey *friends_wiiu_types.GameKey, unknown6 *nex.DateTime) uint32 {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nex.Errors.FPD.InvalidArgument
	}

	client := packet.Sender().(*nex.PRUDPClient)

	senderPID := client.PID().LegacyValue()
	recipientPID := pid

	senderPrincipalInfo, err := utility.GetUserInfoByPID(senderPID)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return nex.Errors.FPD.Unknown
	}

	recipientPrincipalInfo, err := utility.GetUserInfoByPID(recipientPID)
	if err != nil {
		if err == database.ErrPIDNotFound {
			globals.Logger.Errorf("User %d has sent friend request to invalid PID %d", senderPID, pid)
			return nex.Errors.FPD.InvalidPrincipalID // TODO: Not sure if this is the correct error.
		} else {
			globals.Logger.Critical(err.Error())
			return nex.Errors.FPD.Unknown
		}
	}

	currentTimestamp := time.Now()
	expireTimestamp := currentTimestamp.Add(time.Hour * 24 * 29)

	sentTime := nex.NewDateTime(0)
	expireTime := nex.NewDateTime(0)

	sentTime.FromTimestamp(currentTimestamp)
	expireTime.FromTimestamp(expireTimestamp)

	friendRequestID, err := database_wiiu.SaveFriendRequest(senderPID, recipientPID, sentTime.Value(), expireTime.Value(), message)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return nex.Errors.FPD.Unknown
	}

	friendRequest := friends_wiiu_types.NewFriendRequest()

	friendRequest.PrincipalInfo = recipientPrincipalInfo

	friendRequest.Message = friends_wiiu_types.NewFriendRequestMessage()
	friendRequest.Message.FriendRequestID = friendRequestID
	friendRequest.Message.Received = false
	friendRequest.Message.Unknown2 = 1 // replaying from real
	friendRequest.Message.Message = message
	friendRequest.Message.Unknown3 = 0           // replaying from real server
	friendRequest.Message.Unknown4 = ""          // replaying from real server
	friendRequest.Message.GameKey = gameKey      // maybe this is reused?
	friendRequest.Message.Unknown5 = unknown6    // maybe this is reused?
	friendRequest.Message.ExpiresOn = expireTime // no idea why this is set as the sent time
	friendRequest.SentOn = sentTime

	// Why does this exist?? Always empty??
	friendInfo := friends_wiiu_types.NewFriendInfo()

	friendInfo.NNAInfo = friends_wiiu_types.NewNNAInfo()
	friendInfo.NNAInfo.PrincipalBasicInfo = friends_wiiu_types.NewPrincipalBasicInfo()
	friendInfo.NNAInfo.PrincipalBasicInfo.PID = 0
	friendInfo.NNAInfo.PrincipalBasicInfo.NNID = ""
	friendInfo.NNAInfo.PrincipalBasicInfo.Mii = friends_wiiu_types.NewMiiV2()
	friendInfo.NNAInfo.PrincipalBasicInfo.Mii.Name = ""
	friendInfo.NNAInfo.PrincipalBasicInfo.Mii.Unknown1 = 0
	friendInfo.NNAInfo.PrincipalBasicInfo.Mii.Unknown2 = 0
	friendInfo.NNAInfo.PrincipalBasicInfo.Mii.MiiData = []byte{}
	friendInfo.NNAInfo.PrincipalBasicInfo.Mii.Datetime = nex.NewDateTime(0)
	friendInfo.NNAInfo.PrincipalBasicInfo.Unknown = 0
	friendInfo.NNAInfo.Unknown1 = 0
	friendInfo.NNAInfo.Unknown2 = 0

	friendInfo.Presence = friends_wiiu_types.NewNintendoPresenceV2()
	friendInfo.Presence.ChangedFlags = 0
	friendInfo.Presence.Online = false
	friendInfo.Presence.GameKey = gameKey // maybe this is reused?
	friendInfo.Presence.Unknown1 = 0
	friendInfo.Presence.Message = ""
	friendInfo.Presence.Unknown2 = 0
	friendInfo.Presence.Unknown3 = 0
	friendInfo.Presence.GameServerID = 0
	friendInfo.Presence.Unknown4 = 0
	friendInfo.Presence.PID = 0
	friendInfo.Presence.GatheringID = 0
	friendInfo.Presence.ApplicationData = []byte{0x00}
	friendInfo.Presence.Unknown5 = 0
	friendInfo.Presence.Unknown6 = 0
	friendInfo.Presence.Unknown7 = 0

	friendInfo.Status = friends_wiiu_types.NewComment()
	friendInfo.Status.Unknown = 0
	friendInfo.Status.Contents = ""
	friendInfo.Status.LastChanged = nex.NewDateTime(0)

	friendInfo.BecameFriend = nex.NewDateTime(0)
	friendInfo.LastOnline = nex.NewDateTime(0)
	friendInfo.Unknown = 0

	recipientClient := globals.ConnectedUsers[recipientPID]

	if recipientClient != nil {

		friendRequestNotificationData := friends_wiiu_types.NewFriendRequest()

		friendRequestNotificationData.PrincipalInfo = senderPrincipalInfo

		friendRequestNotificationData.Message = friends_wiiu_types.NewFriendRequestMessage()
		friendRequestNotificationData.Message.FriendRequestID = friendRequestID
		friendRequestNotificationData.Message.Received = false
		friendRequestNotificationData.Message.Unknown2 = 1 // replaying from real
		friendRequestNotificationData.Message.Message = message
		friendRequestNotificationData.Message.Unknown3 = 0           // replaying from real server
		friendRequestNotificationData.Message.Unknown4 = ""          // replaying from real server
		friendRequestNotificationData.Message.GameKey = gameKey      // maybe this is reused?
		friendRequestNotificationData.Message.Unknown5 = unknown6    // maybe this is reused?
		friendRequestNotificationData.Message.ExpiresOn = expireTime // no idea why this is set as the sent time
		friendRequestNotificationData.SentOn = sentTime

		go notifications_wiiu.SendFriendRequest(recipientClient.Client, friendRequestNotificationData)
	}

	rmcResponseStream := nex.NewStreamOut(globals.SecureServer)

	rmcResponseStream.WriteStructure(friendRequest)
	rmcResponseStream.WriteStructure(friendInfo)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(rmcResponseBody)
	rmcResponse.ProtocolID = friends_wiiu.ProtocolID
	rmcResponse.MethodID = friends_wiiu.MethodAddFriendRequest
	rmcResponse.CallID = callID

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewPRUDPPacketV0(client, nil)

	responsePacket.SetType(nex.DataPacket)
	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)
	responsePacket.SetSourceStreamType(packet.(nex.PRUDPPacketInterface).DestinationStreamType())
	responsePacket.SetSourcePort(packet.(nex.PRUDPPacketInterface).DestinationPort())
	responsePacket.SetDestinationStreamType(packet.(nex.PRUDPPacketInterface).SourceStreamType())
	responsePacket.SetDestinationPort(packet.(nex.PRUDPPacketInterface).SourcePort())
	responsePacket.SetPayload(rmcResponseBytes)

	globals.SecureServer.Send(responsePacket)

	return 0
}
