package main

import (
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func acceptFriendRequest(err error, client *nex.Client, callID uint32, id uint64) {
	friendInfo := acceptFriendshipAndReturnFriendInfo(id)

	friendPID := friendInfo.NNAInfo.PrincipalBasicInfo.PID
	connectedUser := connectedUsers[friendPID]

	if connectedUser != nil {
		senderPID := client.PID()
		senderConnectedUser := connectedUsers[senderPID]

		senderFriendInfo := nexproto.NewFriendInfo()

		senderFriendInfo.NNAInfo = senderConnectedUser.NNAInfo
		senderFriendInfo.Presence = senderConnectedUser.Presence
		senderFriendInfo.Status = getUserComment(senderPID)
		senderFriendInfo.BecameFriend = friendInfo.BecameFriend
		senderFriendInfo.LastOnline = friendInfo.LastOnline // TODO: Change this
		senderFriendInfo.Unknown = 0

		go sendFriendRequestAcceptedNotification(connectedUser.Client, senderFriendInfo)
	}

	rmcResponseStream := nex.NewStreamOut(nexServer)

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

	nexServer.Send(responsePacket)
}

func sendFriendRequestAcceptedNotification(client *nex.Client, friendInfo *nexproto.FriendInfo) {
	eventObject := nexproto.NewNintendoNotificationEvent()
	eventObject.Type = 30
	eventObject.SenderPID = friendInfo.NNAInfo.PrincipalBasicInfo.PID
	eventObject.DataHolder = nex.NewDataHolder()
	eventObject.DataHolder.SetTypeName("FriendInfo")
	eventObject.DataHolder.SetObjectData(friendInfo)

	stream := nex.NewStreamOut(nexServer)
	eventObjectBytes := eventObject.Bytes(stream)

	rmcRequest := nex.NewRMCRequest()
	rmcRequest.SetProtocolID(nexproto.NintendoNotificationsProtocolID)
	rmcRequest.SetCallID(3810693103)
	rmcRequest.SetMethodID(nexproto.NintendoNotificationsMethodProcessNintendoNotificationEvent1)
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

	nexServer.Send(requestPacket)
}
