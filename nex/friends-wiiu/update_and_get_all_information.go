package nex_friends_wiiu

import (
	"os"

	"github.com/PretendoNetwork/friends/database"
	database_wiiu "github.com/PretendoNetwork/friends/database/wiiu"
	"github.com/PretendoNetwork/friends/globals"
	notifications_wiiu "github.com/PretendoNetwork/friends/notifications/wiiu"
	"github.com/PretendoNetwork/friends/types"
	nex "github.com/PretendoNetwork/nex-go"
	friends_wiiu "github.com/PretendoNetwork/nex-protocols-go/friends-wiiu"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/friends-wiiu/types"
)

func UpdateAndGetAllInformation(err error, client *nex.Client, callID uint32, nnaInfo *friends_wiiu_types.NNAInfo, presence *friends_wiiu_types.NintendoPresenceV2, birthday *nex.DateTime) uint32 {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nex.Errors.FPD.InvalidArgument
	}

	// Get user information
	pid := client.PID()

	if globals.ConnectedUsers[pid] == nil {
		// TODO - Figure out why this is getting removed
		connectedUser := types.NewConnectedUser()
		connectedUser.PID = pid
		connectedUser.Platform = types.WUP
		connectedUser.Client = client

		globals.ConnectedUsers[pid] = connectedUser
	}

	globals.ConnectedUsers[pid].NNAInfo = nnaInfo
	globals.ConnectedUsers[pid].PresenceV2 = presence

	principalPreference, err := database_wiiu.GetUserPrincipalPreference(pid)
	if err != nil {
		if err == database.ErrPIDNotFound {
			return nex.Errors.FPD.InvalidPrincipalID
		} else {
			globals.Logger.Critical(err.Error())
			return nex.Errors.FPD.Unknown
		}
	}

	comment, err := database_wiiu.GetUserComment(pid)
	if err != nil {
		if err == database.ErrPIDNotFound {
			return nex.Errors.FPD.InvalidPrincipalID
		} else {
			globals.Logger.Critical(err.Error())
			return nex.Errors.FPD.Unknown
		}
	}

	friendList, err := database_wiiu.GetUserFriendList(pid)
	if err != nil && err != database.ErrEmptyList {
		globals.Logger.Critical(err.Error())
		return nex.Errors.FPD.Unknown
	}

	friendRequestsOut, err := database_wiiu.GetUserFriendRequestsOut(pid)
	if err != nil && err != database.ErrEmptyList {
		globals.Logger.Critical(err.Error())
		return nex.Errors.FPD.Unknown
	}

	friendRequestsIn, err := database_wiiu.GetUserFriendRequestsIn(pid)
	if err != nil && err != database.ErrEmptyList {
		globals.Logger.Critical(err.Error())
		return nex.Errors.FPD.Unknown
	}

	blockList, err := database_wiiu.GetUserBlockList(pid)
	if err != nil && err != database.ErrBlacklistNotFound {
		globals.Logger.Critical(err.Error())
		return nex.Errors.FPD.Unknown
	}

	notifications := database_wiiu.GetUserNotifications(pid)

	// Update user information

	presence.Online = true // Force online status. I have no idea why this is always false
	presence.PID = pid     // WHY IS THIS SET TO 0 BY DEFAULT??

	notifications_wiiu.SendPresenceUpdate(presence)

	if os.Getenv("PN_FRIENDS_CONFIG_ENABLE_BELLA") == "true" {
		bella := friends_wiiu_types.NewFriendInfo()

		bella.NNAInfo = friends_wiiu_types.NewNNAInfo()
		bella.Presence = friends_wiiu_types.NewNintendoPresenceV2()
		bella.Status = friends_wiiu_types.NewComment()
		bella.BecameFriend = nex.NewDateTime(0)
		bella.LastOnline = nex.NewDateTime(0)
		bella.Unknown = 0

		bella.NNAInfo.PrincipalBasicInfo = friends_wiiu_types.NewPrincipalBasicInfo()
		bella.NNAInfo.Unknown1 = 0
		bella.NNAInfo.Unknown2 = 0

		bella.NNAInfo.PrincipalBasicInfo.PID = 1743126339
		bella.NNAInfo.PrincipalBasicInfo.NNID = "bells1998"
		bella.NNAInfo.PrincipalBasicInfo.Mii = friends_wiiu_types.NewMiiV2()
		bella.NNAInfo.PrincipalBasicInfo.Unknown = 0

		bella.NNAInfo.PrincipalBasicInfo.Mii.Name = "bella"
		bella.NNAInfo.PrincipalBasicInfo.Mii.Unknown1 = 0
		bella.NNAInfo.PrincipalBasicInfo.Mii.Unknown2 = 0
		bella.NNAInfo.PrincipalBasicInfo.Mii.MiiData = []byte{
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
		bella.Presence.GameKey = friends_wiiu_types.NewGameKey()
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

	rmcResponseStream := nex.NewStreamOut(globals.SecureServer)

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
	rmcResponse := nex.NewRMCResponse(friends_wiiu.ProtocolID, callID)
	rmcResponse.SetSuccess(friends_wiiu.MethodUpdateAndGetAllInformation, rmcResponseBody)

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
