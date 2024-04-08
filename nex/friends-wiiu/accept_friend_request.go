package nex_friends_wiiu

import (
	"github.com/PretendoNetwork/friends/database"
	database_wiiu "github.com/PretendoNetwork/friends/database/wiiu"
	"github.com/PretendoNetwork/friends/globals"
	notifications_wiiu "github.com/PretendoNetwork/friends/notifications/wiiu"
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	friends_wiiu "github.com/PretendoNetwork/nex-protocols-go/v2/friends-wiiu"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/v2/friends-wiiu/types"
)

func AcceptFriendRequest(err error, packet nex.PacketInterface, callID uint32, id *types.PrimitiveU64) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.FPD.InvalidArgument, "") // TODO - Add error message
	}

	connection := packet.Sender().(*nex.PRUDPConnection)

	friendInfo, err := database_wiiu.AcceptFriendRequestAndReturnFriendInfo(id.Value)
	if err != nil {
		if err == database.ErrFriendRequestNotFound {
			return nil, nex.NewError(nex.ResultCodes.FPD.InvalidMessageID, "") // TODO - Add error message
		} else {
			globals.Logger.Critical(err.Error())
			return nil, nex.NewError(nex.ResultCodes.FPD.Unknown, "") // TODO - Add error message
		}
	}

	friendPID := friendInfo.NNAInfo.PrincipalBasicInfo.PID.LegacyValue()
	connectedUser := globals.ConnectedUsers[friendPID]

	if connectedUser != nil {
		senderPID := connection.PID().LegacyValue()
		senderConnectedUser := globals.ConnectedUsers[senderPID]

		senderFriendInfo := friends_wiiu_types.NewFriendInfo()

		senderFriendInfo.NNAInfo = senderConnectedUser.NNAInfo
		senderFriendInfo.Presence = senderConnectedUser.PresenceV2
		status, err := database_wiiu.GetUserComment(senderPID)
		if err != nil {
			globals.Logger.Critical(err.Error())
			senderFriendInfo.Status = friends_wiiu_types.NewComment()
			senderFriendInfo.Status.LastChanged = types.NewDateTime(0)
		} else {
			senderFriendInfo.Status = status
		}

		senderFriendInfo.BecameFriend = friendInfo.BecameFriend
		senderFriendInfo.LastOnline = friendInfo.LastOnline // TODO - Change this
		senderFriendInfo.Unknown = types.NewPrimitiveU64(0)

		go notifications_wiiu.SendFriendRequestAccepted(connectedUser.Connection, senderFriendInfo)
	}

	rmcResponseStream := nex.NewByteStreamOut(globals.SecureEndpoint.LibraryVersions(), globals.SecureEndpoint.ByteStreamSettings())

	friendInfo.WriteTo(rmcResponseStream)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(globals.SecureEndpoint, rmcResponseBody)
	rmcResponse.ProtocolID = friends_wiiu.ProtocolID
	rmcResponse.MethodID = friends_wiiu.MethodAcceptFriendRequest
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
