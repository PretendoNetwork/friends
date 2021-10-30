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

	// Get user information
	pid := client.PID()

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

	/*
		//List<FriendInfo>
		friendList := make([]nex.StructureInterface, 0, 1)

		friend := nexproto.NewFriendInfo()

		friend.NNAInfo = nexproto.NewNNAInfo()
		friend.Presence = nexproto.NewNintendoPresenceV2()
		friend.Status = nexproto.NewComment()
		friend.BecameFriend = nex.NewDateTime(0)
		friend.LastOnline = nex.NewDateTime(0)
		friend.Unknown = 0

		friend.NNAInfo.PrincipalBasicInfo = nexproto.NewPrincipalBasicInfo()
		friend.NNAInfo.Unknown1 = 0
		friend.NNAInfo.Unknown2 = 0

		friend.NNAInfo.PrincipalBasicInfo.PID = 1743126339
		friend.NNAInfo.PrincipalBasicInfo.NNID = "bells1998"
		friend.NNAInfo.PrincipalBasicInfo.Mii = nexproto.NewMiiV2()
		friend.NNAInfo.PrincipalBasicInfo.Unknown = 0

		friend.NNAInfo.PrincipalBasicInfo.Mii.Name = "bella"
		friend.NNAInfo.PrincipalBasicInfo.Mii.Unknown1 = 0
		friend.NNAInfo.PrincipalBasicInfo.Mii.Unknown2 = 0
		friend.NNAInfo.PrincipalBasicInfo.Mii.Data = []byte{
			0x03, 0x00, 0x00, 0x40, 0xE9, 0x55, 0xA2, 0x09,
			0xE7, 0xC7, 0x41, 0x82, 0xD9, 0x7D, 0x0B, 0x2D,
			0x03, 0xB3, 0xB8, 0x8D, 0x27, 0xD9, 0x00, 0x00,
			0x01, 0x40, 0x62, 0x00, 0x65, 0x00, 0x6C, 0x00,
			0x6C, 0x00, 0x61, 0x00, 0x00, 0x00, 0x45, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x40,
			0x12, 0x00, 0x81, 0x01, 0x04, 0x68, 0x43, 0x18,
			0x20, 0x34, 0x46, 0x14, 0x81, 0x12, 0x17, 0x68,
			0x0D, 0x00, 0x00, 0x29, 0x03, 0x52, 0x48, 0x50,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFE, 0x86,
		}
		friend.NNAInfo.PrincipalBasicInfo.Mii.Datetime = nex.NewDateTime(0)

		friend.Presence.ChangedFlags = 0x1EE
		friend.Presence.Online = true
		friend.Presence.GameKey = nexproto.NewGameKey()
		friend.Presence.Unknown1 = 0
		friend.Presence.Message = "Testing"
		//friend.Presence.Unknown2 = 2
		friend.Presence.Unknown2 = 0
		//friend.Presence.Unknown3 = 2
		friend.Presence.Unknown3 = 0
		//friend.Presence.GameServerID = 0x1010EB00
		friend.Presence.GameServerID = 0
		//friend.Presence.Unknown4 = 3
		friend.Presence.Unknown4 = 0
		friend.Presence.PID = 1743126339
		//friend.Presence.GatheringID = 1743126339 // test fake ID
		friend.Presence.GatheringID = 0
		//friend.Presence.ApplicationData, _ = hex.DecodeString("0000200300000000000000001843ffe567000000")
		friend.Presence.ApplicationData = []byte{0x0}
		friend.Presence.Unknown5 = 0
		friend.Presence.Unknown6 = 0
		friend.Presence.Unknown7 = 0

		//friend.Presence.GameKey.TitleID = 0x000500001010EC00
		friend.Presence.GameKey.TitleID = 0
		//friend.Presence.GameKey.TitleVersion = 64
		friend.Presence.GameKey.TitleVersion = 0

		friend.Status.Unknown = 0
		friend.Status.Contents = "test"
		friend.Status.LastChanged = nex.NewDateTime(0)

		friendList = append(friendList, friend)
	*/

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
