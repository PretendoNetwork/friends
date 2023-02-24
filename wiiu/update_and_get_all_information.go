package friends_wiiu

import (
	"os"

	database_wiiu "github.com/PretendoNetwork/friends-secure/database/wiiu"
	"github.com/PretendoNetwork/friends-secure/globals"
	notifications_wiiu "github.com/PretendoNetwork/friends-secure/notifications/wiiu"
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func UpdateAndGetAllInformation(err error, client *nex.Client, callID uint32, nnaInfo *nexproto.NNAInfo, presence *nexproto.NintendoPresenceV2, birthday *nex.DateTime) {

	if err != nil {
		// TODO: Handle error
		globals.Logger.Critical(err.Error())
	}

	// Update user information

	presence.Online = true      // Force online status. I have no idea why this is always false
	presence.PID = client.PID() // WHY IS THIS SET TO 0 BY DEFAULT??

	notifications_wiiu.SendPresenceUpdate(presence)

	// Get user information
	pid := client.PID()

	globals.ConnectedUsers[pid].NNAInfo = nnaInfo
	globals.ConnectedUsers[pid].PresenceV2 = presence

	principalPreference := database_wiiu.GetUserPrincipalPreference(pid)
	comment := database_wiiu.GetUserComment(pid)
	friendList := database_wiiu.GetUserFriendList(pid)
	friendRequestsOut := database_wiiu.GetUserFriendRequestsOut(pid)
	friendRequestsIn := database_wiiu.GetUserFriendRequestsIn(pid)
	blockList := database_wiiu.GetUserBlockList(pid)
	notifications := database_wiiu.GetUserNotifications(pid)

	if os.Getenv("ENABLE_BELLA") == "true" {
		bella := nexproto.NewFriendInfo()

		bella.NNAInfo = nexproto.NewNNAInfo()
		bella.Presence = nexproto.NewNintendoPresenceV2()
		bella.Status = nexproto.NewComment()
		bella.BecameFriend = nex.NewDateTime(0)
		bella.LastOnline = nex.NewDateTime(0)
		bella.Unknown = 0

		bella.NNAInfo.PrincipalBasicInfo = nexproto.NewPrincipalBasicInfo()
		bella.NNAInfo.Unknown1 = 0
		bella.NNAInfo.Unknown2 = 0

		bella.NNAInfo.PrincipalBasicInfo.PID = 1743126339
		bella.NNAInfo.PrincipalBasicInfo.NNID = "bells1998"
		bella.NNAInfo.PrincipalBasicInfo.Mii = nexproto.NewMiiV2()
		bella.NNAInfo.PrincipalBasicInfo.Unknown = 0

		bella.NNAInfo.PrincipalBasicInfo.Mii.Name = "bella"
		bella.NNAInfo.PrincipalBasicInfo.Mii.Unknown1 = 0
		bella.NNAInfo.PrincipalBasicInfo.Mii.Unknown2 = 0
		bella.NNAInfo.PrincipalBasicInfo.Mii.Data = []byte{
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
		bella.NNAInfo.PrincipalBasicInfo.Mii.Datetime = nex.NewDateTime(0)

		bella.Presence.ChangedFlags = 0x1EE
		bella.Presence.Online = true
		bella.Presence.GameKey = nexproto.NewGameKey()
		bella.Presence.Unknown1 = 0
		bella.Presence.Message = "Testing"
		//bella.Presence.Unknown2 = 2
		bella.Presence.Unknown2 = 0
		//bella.Presence.Unknown3 = 2
		bella.Presence.Unknown3 = 0
		//bella.Presence.GameServerID = 0x1010EB00
		bella.Presence.GameServerID = 0
		//bella.Presence.Unknown4 = 3
		bella.Presence.Unknown4 = 0
		bella.Presence.PID = 1743126339
		//bella.Presence.GatheringID = 1743126339 // test fake ID
		bella.Presence.GatheringID = 0
		//bella.Presence.ApplicationData, _ = hex.DecodeString("0000200300000000000000001843ffe567000000")
		bella.Presence.ApplicationData = []byte{0x0}
		bella.Presence.Unknown5 = 0
		bella.Presence.Unknown6 = 0
		bella.Presence.Unknown7 = 0

		//bella.Presence.GameKey.TitleID = 0x000500001010EC00
		bella.Presence.GameKey.TitleID = 0
		//bella.Presence.GameKey.TitleVersion = 64
		bella.Presence.GameKey.TitleVersion = 0

		bella.Status.Unknown = 0
		bella.Status.Contents = "test"
		bella.Status.LastChanged = nex.NewDateTime(0)

		friendList = append(friendList, bella)
	}

	rmcResponseStream := nex.NewStreamOut(globals.NEXServer)

	rmcResponseStream.WriteStructure(principalPreference)
	rmcResponseStream.WriteStructure(comment)
	rmcResponseStream.WriteListStructure(friendList)
	rmcResponseStream.WriteListStructure(friendRequestsOut)
	rmcResponseStream.WriteListStructure(friendRequestsIn)
	rmcResponseStream.WriteListStructure(blockList)
	rmcResponseStream.WriteBool(false) // Unknown
	rmcResponseStream.WriteListStructure(notifications)

	//Unknown Bool
	rmcResponseStream.WriteUInt8(0)

	rmcResponseBody := rmcResponseStream.Bytes()

	// Build response packet
	rmcResponse := nex.NewRMCResponse(nexproto.FriendsWiiUProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.FriendsWiiUMethodUpdateAndGetAllInformation, rmcResponseBody)

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
