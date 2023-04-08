package nex_friends_wiiu

import (
	"fmt"
	"time"

	database_wiiu "github.com/PretendoNetwork/friends-secure/database/wiiu"
	"github.com/PretendoNetwork/friends-secure/globals"
	notifications_wiiu "github.com/PretendoNetwork/friends-secure/notifications/wiiu"
	nex "github.com/PretendoNetwork/nex-go"
	friends_wiiu "github.com/PretendoNetwork/nex-protocols-go/friends/wiiu"
)

func AddFriendRequest(err error, client *nex.Client, callID uint32, pid uint32, unknown2 uint8, message string, unknown4 uint8, unknown5 string, gameKey *friends_wiiu.GameKey, unknown6 *nex.DateTime) {
	senderPID := client.PID()
	recipientPID := pid

	recipient := database_wiiu.GetUserInfoByPID(recipientPID)
	if recipient == nil {
		globals.Logger.Error(fmt.Sprintf("User %d has sent friend request to invalid PID %d", senderPID, pid))

		rmcResponse := nex.NewRMCResponse(friends_wiiu.ProtocolID, callID)
		rmcResponse.SetError(nex.Errors.FPD.InvalidPrincipalID) // TODO - Is this the right error?

		rmcResponseBytes := rmcResponse.Bytes()

		responsePacket, _ := nex.NewPacketV0(client, nil)

		responsePacket.SetVersion(0)
		responsePacket.SetSource(0xA1)
		responsePacket.SetDestination(0xAF)
		responsePacket.SetType(nex.DataPacket)
		responsePacket.SetPayload(rmcResponseBytes)

		responsePacket.AddFlag(nex.FlagNeedsAck)
		responsePacket.AddFlag(nex.FlagReliable)

		globals.NEXServer.Send(responsePacket)

		return
	}

	currentTimestamp := time.Now()
	expireTimestamp := currentTimestamp.Add(time.Hour * 24 * 29)

	sentTime := nex.NewDateTime(0)
	expireTime := nex.NewDateTime(0)

	sentTime.FromTimestamp(currentTimestamp)
	expireTime.FromTimestamp(expireTimestamp)

	friendRequestID := database_wiiu.SaveFriendRequest(senderPID, recipientPID, sentTime.Value(), expireTime.Value(), message)

	friendRequest := friends_wiiu.NewFriendRequest()

	friendRequest.PrincipalInfo = database_wiiu.GetUserInfoByPID(recipientPID)

	friendRequest.Message = friends_wiiu.NewFriendRequestMessage()
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
	friendInfo := friends_wiiu.NewFriendInfo()

	friendInfo.NNAInfo = friends_wiiu.NewNNAInfo()
	friendInfo.NNAInfo.PrincipalBasicInfo = friends_wiiu.NewPrincipalBasicInfo()
	friendInfo.NNAInfo.PrincipalBasicInfo.PID = 0
	friendInfo.NNAInfo.PrincipalBasicInfo.NNID = ""
	friendInfo.NNAInfo.PrincipalBasicInfo.Mii = friends_wiiu.NewMiiV2()
	friendInfo.NNAInfo.PrincipalBasicInfo.Mii.Name = ""
	friendInfo.NNAInfo.PrincipalBasicInfo.Mii.Unknown1 = 0
	friendInfo.NNAInfo.PrincipalBasicInfo.Mii.Unknown2 = 0
	friendInfo.NNAInfo.PrincipalBasicInfo.Mii.Data = []byte{}
	friendInfo.NNAInfo.PrincipalBasicInfo.Mii.Datetime = nex.NewDateTime(0)
	friendInfo.NNAInfo.PrincipalBasicInfo.Unknown = 0
	friendInfo.NNAInfo.Unknown1 = 0
	friendInfo.NNAInfo.Unknown2 = 0

	friendInfo.Presence = friends_wiiu.NewNintendoPresenceV2()
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

	friendInfo.Status = friends_wiiu.NewComment()
	friendInfo.Status.Unknown = 0
	friendInfo.Status.Contents = ""
	friendInfo.Status.LastChanged = nex.NewDateTime(0)

	friendInfo.BecameFriend = nex.NewDateTime(0)
	friendInfo.LastOnline = nex.NewDateTime(0)
	friendInfo.Unknown = 0

	recipientClient := client.Server().FindClientFromPID(recipientPID)

	if recipientClient != nil {

		friendRequestNotificationData := friends_wiiu.NewFriendRequest()

		friendRequestNotificationData.PrincipalInfo = database_wiiu.GetUserInfoByPID(senderPID)

		friendRequestNotificationData.Message = friends_wiiu.NewFriendRequestMessage()
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

		go notifications_wiiu.SendFriendRequest(recipientClient, friendRequestNotificationData)
	}

	rmcResponseStream := nex.NewStreamOut(globals.NEXServer)

	rmcResponseStream.WriteStructure(friendRequest)
	rmcResponseStream.WriteStructure(friendInfo)

	rmcResponseBody := rmcResponseStream.Bytes()

	// Build response packet
	rmcResponse := nex.NewRMCResponse(friends_wiiu.ProtocolID, callID)
	rmcResponse.SetSuccess(friends_wiiu.MethodAddFriendRequest, rmcResponseBody)

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewPacketV0(client, nil)

	responsePacket.SetVersion(0)
	responsePacket.SetSource(0xA1)
	responsePacket.SetDestination(0xAF)
	responsePacket.SetType(nex.DataPacket)
	responsePacket.SetPayload(rmcResponseBytes)

	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)

	globals.NEXServer.Send(responsePacket)
}
