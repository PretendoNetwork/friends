package friends_wiiu

import (
	"encoding/base64"
	"math/rand"
	"time"

	"github.com/PretendoNetwork/friends-secure/database"
	"github.com/PretendoNetwork/friends-secure/globals"
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
	"go.mongodb.org/mongo-driver/bson"
)

func AddFriendRequest(err error, client *nex.Client, callID uint32, pid uint32, unknown2 uint8, message string, unknown4 uint8, unknown5 string, gameKey *nexproto.GameKey, unknown6 *nex.DateTime) {
	rand.Seed(time.Now().UnixNano())
	nodeID := rand.Intn(len(globals.SnowflakeNodes))

	snowflakeNode := globals.SnowflakeNodes[nodeID]

	friendRequestID := uint64(snowflakeNode.Generate().Int64())

	currentTimestamp := time.Now()
	expireTimestamp := currentTimestamp.Add(time.Hour * 24 * 29)

	sentTime := nex.NewDateTime(0)
	expireTime := nex.NewDateTime(0)
	sentTime.FromTimestamp(currentTimestamp)
	expireTime.FromTimestamp(expireTimestamp)

	senderPID := client.PID()
	recipientPID := pid

	recipientUserInforation := database.GetUserInfoByPID(recipientPID)

	friendRequest := nexproto.NewFriendRequest()

	friendRequest.PrincipalInfo = nexproto.NewPrincipalBasicInfo()
	friendRequest.PrincipalInfo.PID = recipientPID
	friendRequest.PrincipalInfo.NNID = recipientUserInforation["username"].(string)
	friendRequest.PrincipalInfo.Mii = nexproto.NewMiiV2()
	friendRequest.PrincipalInfo.Unknown = 2 // replaying from real server

	encodedMiiData := recipientUserInforation["mii"].(bson.M)["data"].(string)
	decodedMiiData, _ := base64.StdEncoding.DecodeString(encodedMiiData)

	friendRequest.PrincipalInfo.Mii.Name = recipientUserInforation["mii"].(bson.M)["name"].(string)
	friendRequest.PrincipalInfo.Mii.Unknown1 = 0 // replaying from real server
	friendRequest.PrincipalInfo.Mii.Unknown2 = 0 // replaying from real server
	friendRequest.PrincipalInfo.Mii.Data = decodedMiiData
	friendRequest.PrincipalInfo.Mii.Datetime = nex.NewDateTime(0)

	friendRequest.Message = nexproto.NewFriendRequestMessage()
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
	friendInfo := nexproto.NewFriendInfo()

	friendInfo.NNAInfo = nexproto.NewNNAInfo()
	friendInfo.NNAInfo.PrincipalBasicInfo = nexproto.NewPrincipalBasicInfo()
	friendInfo.NNAInfo.PrincipalBasicInfo.PID = 0
	friendInfo.NNAInfo.PrincipalBasicInfo.NNID = ""
	friendInfo.NNAInfo.PrincipalBasicInfo.Mii = nexproto.NewMiiV2()
	friendInfo.NNAInfo.PrincipalBasicInfo.Mii.Name = ""
	friendInfo.NNAInfo.PrincipalBasicInfo.Mii.Unknown1 = 0
	friendInfo.NNAInfo.PrincipalBasicInfo.Mii.Unknown2 = 0
	friendInfo.NNAInfo.PrincipalBasicInfo.Mii.Data = []byte{}
	friendInfo.NNAInfo.PrincipalBasicInfo.Mii.Datetime = nex.NewDateTime(0)
	friendInfo.NNAInfo.PrincipalBasicInfo.Unknown = 0
	friendInfo.NNAInfo.Unknown1 = 0
	friendInfo.NNAInfo.Unknown2 = 0

	friendInfo.Presence = nexproto.NewNintendoPresenceV2()
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

	friendInfo.Status = nexproto.NewComment()
	friendInfo.Status.Unknown = 0
	friendInfo.Status.Contents = ""
	friendInfo.Status.LastChanged = nex.NewDateTime(0)

	friendInfo.BecameFriend = nex.NewDateTime(0)
	friendInfo.LastOnline = nex.NewDateTime(0)
	friendInfo.Unknown = 0

	database.SaveFriendRequest(friendRequestID, senderPID, recipientPID, sentTime.Value(), expireTime.Value(), message)

	recipientClient := client.Server().FindClientFromPID(recipientPID)

	if recipientClient != nil {
		senderUserInforation := database.GetUserInfoByPID(senderPID)

		friendRequestNotificationData := nexproto.NewFriendRequest()

		friendRequestNotificationData.PrincipalInfo = nexproto.NewPrincipalBasicInfo()
		friendRequestNotificationData.PrincipalInfo.PID = senderPID
		friendRequestNotificationData.PrincipalInfo.NNID = senderUserInforation["username"].(string)
		friendRequestNotificationData.PrincipalInfo.Mii = nexproto.NewMiiV2()
		friendRequestNotificationData.PrincipalInfo.Unknown = 2 // replaying from real server

		encodedMiiData := senderUserInforation["mii"].(bson.M)["data"].(string)
		decodedMiiData, _ := base64.StdEncoding.DecodeString(encodedMiiData)

		friendRequestNotificationData.PrincipalInfo.Mii.Name = senderUserInforation["mii"].(bson.M)["name"].(string)
		friendRequestNotificationData.PrincipalInfo.Mii.Unknown1 = 0 // replaying from real server
		friendRequestNotificationData.PrincipalInfo.Mii.Unknown2 = 0 // replaying from real server
		friendRequestNotificationData.PrincipalInfo.Mii.Data = decodedMiiData
		friendRequestNotificationData.PrincipalInfo.Mii.Datetime = nex.NewDateTime(0)

		friendRequestNotificationData.Message = nexproto.NewFriendRequestMessage()
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

		go sendFriendRequestNotification(recipientClient, friendRequestNotificationData)
	}

	rmcResponseStream := nex.NewStreamOut(globals.NEXServer)

	rmcResponseStream.WriteStructure(friendRequest)
	rmcResponseStream.WriteStructure(friendInfo)

	rmcResponseBody := rmcResponseStream.Bytes()

	// Build response packet
	rmcResponse := nex.NewRMCResponse(nexproto.FriendsWiiUProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.FriendsWiiUMethodAddFriendRequest, rmcResponseBody)

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

func sendFriendRequestNotification(client *nex.Client, friendRequestNotificationData *nexproto.FriendRequest) {
	eventObject := nexproto.NewNintendoNotificationEvent()
	eventObject.Type = 27
	eventObject.SenderPID = friendRequestNotificationData.PrincipalInfo.PID
	eventObject.DataHolder = nex.NewDataHolder()
	eventObject.DataHolder.SetTypeName("FriendRequest")
	eventObject.DataHolder.SetObjectData(friendRequestNotificationData)

	stream := nex.NewStreamOut(globals.NEXServer)
	eventObjectBytes := eventObject.Bytes(stream)

	rmcRequest := nex.NewRMCRequest()
	rmcRequest.SetProtocolID(nexproto.NintendoNotificationsProtocolID)
	rmcRequest.SetCallID(3810693103)
	rmcRequest.SetMethodID(nexproto.NintendoNotificationsMethodProcessNintendoNotificationEvent2)
	rmcRequest.SetParameters(eventObjectBytes)

	rmcRequestBytes := rmcRequest.Bytes()

	requestPacket, _ := nex.NewPacketV0(client, nil)

	requestPacket.SetVersion(0)
	requestPacket.SetSource(0xA1)
	requestPacket.SetDestination(0xAF)
	requestPacket.SetType(nex.DataPacket)
	requestPacket.SetPayload(rmcRequestBytes)

	requestPacket.AddFlag(nex.FlagNeedsAck)
	requestPacket.AddFlag(nex.FlagReliable)

	globals.NEXServer.Send(requestPacket)
}
