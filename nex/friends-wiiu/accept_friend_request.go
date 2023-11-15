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

func AcceptFriendRequest(err error, packet nex.PacketInterface, callID uint32, id uint64) (*nex.RMCMessage, uint32) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.Errors.FPD.InvalidArgument
	}

	client := packet.Sender().(*nex.PRUDPClient)

	friendInfo, err := database_wiiu.AcceptFriendRequestAndReturnFriendInfo(id)
	if err != nil {
		if err == database.ErrFriendRequestNotFound {
			return nil, nex.Errors.FPD.InvalidMessageID
		} else {
			globals.Logger.Critical(err.Error())
			return nil, nex.Errors.FPD.Unknown
		}
	}

	friendPID := friendInfo.NNAInfo.PrincipalBasicInfo.PID.LegacyValue()
	connectedUser := globals.ConnectedUsers[friendPID]

	if connectedUser != nil {
		senderPID := client.PID().LegacyValue()
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
		senderFriendInfo.LastOnline = friendInfo.LastOnline // TODO - Change this
		senderFriendInfo.Unknown = 0

		go notifications_wiiu.SendFriendRequestAccepted(connectedUser.Client, senderFriendInfo)
	}

	rmcResponseStream := nex.NewStreamOut(globals.SecureServer)

	rmcResponseStream.WriteStructure(friendInfo)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(rmcResponseBody)
	rmcResponse.ProtocolID = friends_wiiu.ProtocolID
	rmcResponse.MethodID = friends_wiiu.MethodAcceptFriendRequest
	rmcResponse.CallID = callID

	return rmcResponse, 0
}
