package nex_friends_wiiu

import (
	"time"

	"github.com/PretendoNetwork/friends/database"
	database_wiiu "github.com/PretendoNetwork/friends/database/wiiu"
	"github.com/PretendoNetwork/friends/globals"
	notifications_wiiu "github.com/PretendoNetwork/friends/notifications/wiiu"
	"github.com/PretendoNetwork/friends/utility"
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	friends_wiiu "github.com/PretendoNetwork/nex-protocols-go/v2/friends-wiiu"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/v2/friends-wiiu/types"
)

func AddFriendRequest(err error, packet nex.PacketInterface, callID uint32, pid *types.PID, unknown2 *types.PrimitiveU8, message *types.String, unknown4 *types.PrimitiveU8, unknown5 *types.String, gameKey *friends_wiiu_types.GameKey, unknown6 *types.DateTime) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.FPD.InvalidArgument, "") // TODO - Add error message
	}

	connection := packet.Sender().(*nex.PRUDPConnection)

	senderPID := connection.PID().LegacyValue()
	recipientPID := pid.LegacyValue()

	senderPrincipalInfo, err := utility.GetUserInfoByPID(senderPID)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return nil, nex.NewError(nex.ResultCodes.FPD.Unknown, "") // TODO - Add error message
	}

	recipientPrincipalInfo, err := utility.GetUserInfoByPID(recipientPID)
	if err != nil {
		if err == database.ErrPIDNotFound {
			// TODO - Not sure if this is the correct error.
			globals.Logger.Errorf("User %d has sent friend request to invalid PID %d", senderPID, pid)
			return nil, nex.NewError(nex.ResultCodes.FPD.InvalidPrincipalID, "") // TODO - Add error message
		} else {
			globals.Logger.Critical(err.Error())
			return nil, nex.NewError(nex.ResultCodes.FPD.Unknown, "") // TODO - Add error message
		}
	}

	currentTimestamp := time.Now()
	expireTimestamp := currentTimestamp.Add(time.Hour * 24 * 29)

	sentTime := types.NewDateTime(0)
	expireTime := types.NewDateTime(0)

	sentTime.FromTimestamp(currentTimestamp)
	expireTime.FromTimestamp(expireTimestamp)

	friendRequestID, err := database_wiiu.SaveFriendRequest(senderPID, recipientPID, sentTime.Value(), expireTime.Value(), message.Value)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return nil, nex.NewError(nex.ResultCodes.FPD.Unknown, "") // TODO - Add error message
	}

	friendRequest := friends_wiiu_types.NewFriendRequest()

	friendRequest.PrincipalInfo = recipientPrincipalInfo

	friendRequest.Message = friends_wiiu_types.NewFriendRequestMessage()
	friendRequest.Message.FriendRequestID = types.NewPrimitiveU64(friendRequestID)
	friendRequest.Message.Received = types.NewPrimitiveBool(false)
	friendRequest.Message.Unknown2 = types.NewPrimitiveU8(1) // * Replaying from official server
	friendRequest.Message.Message = message
	friendRequest.Message.Unknown3 = types.NewPrimitiveU8(0) // * Replaying from official server
	friendRequest.Message.Unknown4 = types.NewString("")     // * Replaying from official server
	friendRequest.Message.GameKey = gameKey                  // * Maybe this is reused?
	friendRequest.Message.Unknown5 = unknown6                // * Maybe this is reused?
	friendRequest.Message.ExpiresOn = expireTime             // * No idea why this is set as the sent time
	friendRequest.SentOn = sentTime

	// * Why does this exist?? Always empty??
	friendInfo := friends_wiiu_types.NewFriendInfo()

	friendInfo.NNAInfo = friends_wiiu_types.NewNNAInfo()
	friendInfo.NNAInfo.PrincipalBasicInfo = friends_wiiu_types.NewPrincipalBasicInfo()
	friendInfo.NNAInfo.PrincipalBasicInfo.PID = types.NewPID(0)
	friendInfo.NNAInfo.PrincipalBasicInfo.NNID = types.NewString("")
	friendInfo.NNAInfo.PrincipalBasicInfo.Mii = friends_wiiu_types.NewMiiV2()
	friendInfo.NNAInfo.PrincipalBasicInfo.Mii.Name = types.NewString("")
	friendInfo.NNAInfo.PrincipalBasicInfo.Mii.Unknown1 = types.NewPrimitiveU8(0)
	friendInfo.NNAInfo.PrincipalBasicInfo.Mii.Unknown2 = types.NewPrimitiveU8(0)
	friendInfo.NNAInfo.PrincipalBasicInfo.Mii.MiiData = types.NewBuffer([]byte{})
	friendInfo.NNAInfo.PrincipalBasicInfo.Mii.Datetime = types.NewDateTime(0)
	friendInfo.NNAInfo.PrincipalBasicInfo.Unknown = types.NewPrimitiveU8(0)
	friendInfo.NNAInfo.Unknown1 = types.NewPrimitiveU8(0)
	friendInfo.NNAInfo.Unknown2 = types.NewPrimitiveU8(0)

	friendInfo.Presence = friends_wiiu_types.NewNintendoPresenceV2()
	friendInfo.Presence.ChangedFlags = types.NewPrimitiveU32(0)
	friendInfo.Presence.Online = types.NewPrimitiveBool(false)
	friendInfo.Presence.GameKey = gameKey // * Maybe this is reused?
	friendInfo.Presence.Unknown1 = types.NewPrimitiveU8(0)
	friendInfo.Presence.Message = types.NewString("")
	friendInfo.Presence.Unknown2 = types.NewPrimitiveU32(0)
	friendInfo.Presence.Unknown3 = types.NewPrimitiveU8(0)
	friendInfo.Presence.GameServerID = types.NewPrimitiveU32(0)
	friendInfo.Presence.Unknown4 = types.NewPrimitiveU32(0)
	friendInfo.Presence.PID = types.NewPID(0)
	friendInfo.Presence.GatheringID = types.NewPrimitiveU32(0)
	friendInfo.Presence.ApplicationData = types.NewBuffer([]byte{0x00})
	friendInfo.Presence.Unknown5 = types.NewPrimitiveU8(0)
	friendInfo.Presence.Unknown6 = types.NewPrimitiveU8(0)
	friendInfo.Presence.Unknown7 = types.NewPrimitiveU8(0)

	friendInfo.Status = friends_wiiu_types.NewComment()
	friendInfo.Status.Unknown = types.NewPrimitiveU8(0)
	friendInfo.Status.Contents = types.NewString("")
	friendInfo.Status.LastChanged = types.NewDateTime(0)

	friendInfo.BecameFriend = types.NewDateTime(0)
	friendInfo.LastOnline = types.NewDateTime(0)
	friendInfo.Unknown = types.NewPrimitiveU64(0)

	recipientClient := globals.ConnectedUsers[recipientPID]

	if recipientClient != nil {
		friendRequestNotificationData := friends_wiiu_types.NewFriendRequest()

		friendRequestNotificationData.PrincipalInfo = senderPrincipalInfo
		friendRequestNotificationData.Message = friends_wiiu_types.NewFriendRequestMessage()
		friendRequestNotificationData.Message.FriendRequestID = types.NewPrimitiveU64(friendRequestID)
		friendRequestNotificationData.Message.Received = types.NewPrimitiveBool(false)
		friendRequestNotificationData.Message.Unknown2 = types.NewPrimitiveU8(1) // * Replaying from official server
		friendRequestNotificationData.Message.Message = message
		friendRequestNotificationData.Message.Unknown3 = types.NewPrimitiveU8(0) // * Replaying from server server
		friendRequestNotificationData.Message.Unknown4 = types.NewString("")     // * Replaying from server server
		friendRequestNotificationData.Message.GameKey = gameKey                  // * Maybe this is reused?
		friendRequestNotificationData.Message.Unknown5 = unknown6                // * Maybe this is reused?
		friendRequestNotificationData.Message.ExpiresOn = expireTime             // * No idea why this is set as the sent time
		friendRequestNotificationData.SentOn = sentTime

		go notifications_wiiu.SendFriendRequest(recipientClient.Connection, friendRequestNotificationData)
	}

	rmcResponseStream := nex.NewByteStreamOut(globals.SecureEndpoint.LibraryVersions(), globals.SecureEndpoint.ByteStreamSettings())

	friendRequest.WriteTo(rmcResponseStream)
	friendInfo.WriteTo(rmcResponseStream)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(globals.SecureEndpoint, rmcResponseBody)
	rmcResponse.ProtocolID = friends_wiiu.ProtocolID
	rmcResponse.MethodID = friends_wiiu.MethodAddFriendRequest
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
