package nex_friends_wiiu

import (
	"github.com/PretendoNetwork/friends/database"
	database_wiiu "github.com/PretendoNetwork/friends/database/wiiu"
	"github.com/PretendoNetwork/friends/globals"
	notifications_wiiu "github.com/PretendoNetwork/friends/notifications/wiiu"
	nex "github.com/PretendoNetwork/nex-go"
	friends_wiiu "github.com/PretendoNetwork/nex-protocols-go/friends-wiiu"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/friends-wiiu/types"
)

func AcceptFriendRequest(err error, client *nex.Client, callID uint32, id uint64) uint32 {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nex.Errors.FPD.InvalidArgument
	}

	friendInfo, err := database_wiiu.AcceptFriendRequestAndReturnFriendInfo(id)
	if err != nil {
		if err == database.ErrFriendRequestNotFound {
			return nex.Errors.FPD.InvalidMessageID
		} else {
			globals.Logger.Critical(err.Error())
			return nex.Errors.FPD.Unknown
		}
	}

	friendPID := friendInfo.NNAInfo.PrincipalBasicInfo.PID
	connectedUser := globals.ConnectedUsers[friendPID]

	if connectedUser != nil {
		senderPID := client.PID()
		senderConnectedUser := globals.ConnectedUsers[senderPID]

		senderFriendInfo := friends_wiiu_types.NewFriendInfo()

		senderFriendInfo.NNAInfo = senderConnectedUser.NNAInfo
		senderFriendInfo.Presence = senderConnectedUser.PresenceV2
		status, err := database_wiiu.GetUserComment(senderPID)
		if err != nil {
			globals.Logger.Critical(err.Error())
			senderFriendInfo.Status = friends_wiiu_types.NewComment()
			senderFriendInfo.Status.LastChanged = nex.NewDateTime(0)
		} else {
			senderFriendInfo.Status = status
		}

		senderFriendInfo.BecameFriend = friendInfo.BecameFriend
		senderFriendInfo.LastOnline = friendInfo.LastOnline // TODO: Change this
		senderFriendInfo.Unknown = 0

		go notifications_wiiu.SendFriendRequestAccepted(connectedUser.Client, senderFriendInfo)
	}

	rmcResponseStream := nex.NewStreamOut(globals.SecureServer)

	rmcResponseStream.WriteStructure(friendInfo)

	rmcResponseBody := rmcResponseStream.Bytes()

	// Build response packet
	rmcResponse := nex.NewRMCResponse(friends_wiiu.ProtocolID, callID)
	rmcResponse.SetSuccess(friends_wiiu.MethodAcceptFriendRequest, rmcResponseBody)

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewPacketV0(client, nil)

	responsePacket.SetVersion(0)
	responsePacket.SetSource(0xA1)
	responsePacket.SetDestination(0xAF)
	responsePacket.SetType(nex.DataPacket)
	responsePacket.SetPayload(rmcResponseBytes)

	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)

	globals.SecureServer.Send(responsePacket)

	return 0
}
