package main

import (
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func updateAndGetAllInformation(err error, client *nex.Client, callID uint32, nnaInfo *nexproto.NNAInfo, presence *nexproto.NintendoPresenceV2, birthday *nex.DateTime) {

	if err != nil {
		// TODO: Handle error
		panic(err)
	}

	// Update user information

	updateNNAInfo(nnaInfo)
	updateNintendoPresenceV2(presence)
	go sendUpdatePresenceWiiUNotifications(presence)

	// Get user information
	pid := client.PID()

	connectedUsers[pid].NNAInfo = nnaInfo
	connectedUsers[pid].Presence = presence

	principalPreference := getUserPrincipalPreference(pid)
	comment := getUserComment(pid)
	friendList := getUserFriendList(pid)
	friendRequestsOut := getUserFriendRequestsOut(pid)
	friendRequestsIn := getUserFriendRequestsIn(pid)
	//blockList := getUserBlockList(pid)
	//notifications := getUserNotifications(pid)

	rmcResponseStream := nex.NewStreamOut(nexServer)

	rmcResponseStream.WriteStructure(principalPreference)
	rmcResponseStream.WriteStructure(comment)
	rmcResponseStream.WriteListStructure(friendList)
	rmcResponseStream.WriteListStructure(friendRequestsOut)
	rmcResponseStream.WriteListStructure(friendRequestsIn)
	// End of hard-coded friend

	//List<BlacklistedPrincipal>
	rmcResponseStream.WriteUInt32LE(0)

	//Unknown Bool
	rmcResponseStream.WriteUInt8(0)

	//List<PersistentNotification>
	rmcResponseStream.WriteUInt32LE(0)

	//Unknown Bool
	rmcResponseStream.WriteUInt8(0)

	rmcResponseBody := rmcResponseStream.Bytes()

	// Build response packet
	rmcResponse := nex.NewRMCResponse(nexproto.FriendsProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.FriendsMethodUpdateAndGetAllInformation, rmcResponseBody)

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
