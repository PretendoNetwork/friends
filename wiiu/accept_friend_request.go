package friends_wiiu

import (
	database_wiiu "github.com/PretendoNetwork/friends-secure/database/wiiu"
	"github.com/PretendoNetwork/friends-secure/globals"
	notifications_wiiu "github.com/PretendoNetwork/friends-secure/notifications/wiiu"
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func AcceptFriendRequest(err error, client *nex.Client, callID uint32, id uint64) {
	friendInfo := database_wiiu.AcceptFriendshipAndReturnFriendInfo(id)

	friendPID := friendInfo.NNAInfo.PrincipalBasicInfo.PID
	connectedUser := globals.ConnectedUsers[friendPID]

	if connectedUser != nil {
		senderPID := client.PID()
		senderConnectedUser := globals.ConnectedUsers[senderPID]

		senderFriendInfo := nexproto.NewFriendInfo()

		senderFriendInfo.NNAInfo = senderConnectedUser.NNAInfo
		senderFriendInfo.Presence = senderConnectedUser.PresenceV2
		senderFriendInfo.Status = database_wiiu.GetUserComment(senderPID)
		senderFriendInfo.BecameFriend = friendInfo.BecameFriend
		senderFriendInfo.LastOnline = friendInfo.LastOnline // TODO: Change this
		senderFriendInfo.Unknown = 0

		go notifications_wiiu.SendFriendRequestAccepted(connectedUser.Client, senderFriendInfo)
	}

	rmcResponseStream := nex.NewStreamOut(globals.NEXServer)

	rmcResponseStream.WriteStructure(friendInfo)

	rmcResponseBody := rmcResponseStream.Bytes()

	// Build response packet
	rmcResponse := nex.NewRMCResponse(nexproto.FriendsWiiUProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.FriendsWiiUMethodAcceptFriendRequest, rmcResponseBody)

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
