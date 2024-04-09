package nex_friends_wiiu

import (
	"os"

	"github.com/PretendoNetwork/friends/database"
	database_wiiu "github.com/PretendoNetwork/friends/database/wiiu"
	"github.com/PretendoNetwork/friends/globals"
	notifications_wiiu "github.com/PretendoNetwork/friends/notifications/wiiu"
	friends_types "github.com/PretendoNetwork/friends/types"
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	friends_wiiu "github.com/PretendoNetwork/nex-protocols-go/v2/friends-wiiu"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/v2/friends-wiiu/types"
)

func UpdateAndGetAllInformation(err error, packet nex.PacketInterface, callID uint32, nnaInfo *friends_wiiu_types.NNAInfo, presence *friends_wiiu_types.NintendoPresenceV2, birthday *types.DateTime) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.FPD.InvalidArgument, "") // TODO - Add error message
	}

	connection := packet.Sender().(*nex.PRUDPConnection)

	// * Get user information
	pid := connection.PID().LegacyValue()
	connectedUser, ok := globals.ConnectedUsers.Get(pid)

	if !ok || connectedUser == nil {
		// TODO - Figure out why this is getting removed
		connectedUser := friends_types.NewConnectedUser()
		connectedUser.PID = pid
		connectedUser.Platform = friends_types.WUP
		connectedUser.Connection = connection
		// TODO - Find a clean way to create a NNAInfo?

		globals.ConnectedUsers.Set(pid, connectedUser)
	}

	connectedUser.NNAInfo = nnaInfo
	connectedUser.PresenceV2 = presence

	principalPreference, err := database_wiiu.GetUserPrincipalPreference(pid)
	if err != nil {
		if err == database.ErrPIDNotFound {
			return nil, nex.NewError(nex.ResultCodes.FPD.InvalidPrincipalID, "") // TODO - Add error message
		} else {
			globals.Logger.Critical(err.Error())
			return nil, nex.NewError(nex.ResultCodes.FPD.Unknown, "") // TODO - Add error message
		}
	}

	comment, err := database_wiiu.GetUserComment(pid)
	if err != nil {
		if err == database.ErrPIDNotFound {
			return nil, nex.NewError(nex.ResultCodes.FPD.InvalidPrincipalID, "") // TODO - Add error message
		} else {
			globals.Logger.Critical(err.Error())
			return nil, nex.NewError(nex.ResultCodes.FPD.Unknown, "") // TODO - Add error message
		}
	}

	friendList, err := database_wiiu.GetUserFriendList(pid)
	if err != nil && err != database.ErrEmptyList {
		globals.Logger.Critical(err.Error())
		return nil, nex.NewError(nex.ResultCodes.FPD.Unknown, "") // TODO - Add error message
	}

	friendRequestsOut, err := database_wiiu.GetUserFriendRequestsOut(pid)
	if err != nil && err != database.ErrEmptyList {
		globals.Logger.Critical(err.Error())
		return nil, nex.NewError(nex.ResultCodes.FPD.Unknown, "") // TODO - Add error message
	}

	friendRequestsIn, err := database_wiiu.GetUserFriendRequestsIn(pid)
	if err != nil && err != database.ErrEmptyList {
		globals.Logger.Critical(err.Error())
		return nil, nex.NewError(nex.ResultCodes.FPD.Unknown, "") // TODO - Add error message
	}

	blockList, err := database_wiiu.GetUserBlockList(pid)
	if err != nil && err != database.ErrBlacklistNotFound {
		globals.Logger.Critical(err.Error())
		return nil, nex.NewError(nex.ResultCodes.FPD.Unknown, "") // TODO - Add error message
	}

	notifications := database_wiiu.GetUserNotifications(pid)

	// * Update user information

	presence.Online = types.NewPrimitiveBool(true) // * Force online status. I have no idea why this is always false
	presence.PID = connection.PID()                // * WHY IS THIS SET TO 0 BY DEFAULT??

	notifications_wiiu.SendPresenceUpdate(presence)

	if os.Getenv("PN_FRIENDS_CONFIG_ENABLE_BELLA") == "true" {
		bella := friends_wiiu_types.NewFriendInfo()

		bella.NNAInfo = friends_wiiu_types.NewNNAInfo()
		bella.Presence = friends_wiiu_types.NewNintendoPresenceV2()
		bella.Status = friends_wiiu_types.NewComment()
		bella.BecameFriend = types.NewDateTime(0)
		bella.LastOnline = types.NewDateTime(0)
		bella.Unknown = types.NewPrimitiveU64(0)

		bella.NNAInfo.PrincipalBasicInfo = friends_wiiu_types.NewPrincipalBasicInfo()
		bella.NNAInfo.Unknown1 = types.NewPrimitiveU8(0)
		bella.NNAInfo.Unknown2 = types.NewPrimitiveU8(0)

		bella.NNAInfo.PrincipalBasicInfo.PID = types.NewPID(1743126339)
		bella.NNAInfo.PrincipalBasicInfo.NNID = types.NewString("bells1998")
		bella.NNAInfo.PrincipalBasicInfo.Mii = friends_wiiu_types.NewMiiV2()
		bella.NNAInfo.PrincipalBasicInfo.Unknown = types.NewPrimitiveU8(0)

		bella.NNAInfo.PrincipalBasicInfo.Mii.Name = types.NewString("bella")
		bella.NNAInfo.PrincipalBasicInfo.Mii.Unknown1 = types.NewPrimitiveU8(0)
		bella.NNAInfo.PrincipalBasicInfo.Mii.Unknown2 = types.NewPrimitiveU8(0)
		bella.NNAInfo.PrincipalBasicInfo.Mii.MiiData = types.NewBuffer([]byte{
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
		})
		bella.NNAInfo.PrincipalBasicInfo.Mii.Datetime = types.NewDateTime(0)

		bella.Presence.ChangedFlags = types.NewPrimitiveU32(0x1EE)
		bella.Presence.Online = types.NewPrimitiveBool(true)
		bella.Presence.GameKey = friends_wiiu_types.NewGameKey()
		bella.Presence.Unknown1 = types.NewPrimitiveU8(0)
		bella.Presence.Message = types.NewString("Testing")
		//bella.Presence.Unknown2 = 2
		bella.Presence.Unknown2 = types.NewPrimitiveU32(0)
		//bella.Presence.Unknown3 = 2
		bella.Presence.Unknown3 = types.NewPrimitiveU8(0)
		//bella.Presence.GameServerID = 0x1010EB00
		bella.Presence.GameServerID = types.NewPrimitiveU32(0)
		//bella.Presence.Unknown4 = 3
		bella.Presence.Unknown4 = types.NewPrimitiveU32(0)
		bella.Presence.PID = types.NewPID(1743126339)
		//bella.Presence.GatheringID = 1743126339 // test fake ID
		bella.Presence.GatheringID = types.NewPrimitiveU32(0)
		//bella.Presence.ApplicationData, _ = hex.DecodeString("0000200300000000000000001843ffe567000000")
		bella.Presence.ApplicationData = types.NewBuffer([]byte{0x0})
		bella.Presence.Unknown5 = types.NewPrimitiveU8(0)
		bella.Presence.Unknown6 = types.NewPrimitiveU8(0)
		bella.Presence.Unknown7 = types.NewPrimitiveU8(0)

		//bella.Presence.GameKey.TitleID = 0x000500001010EC00
		bella.Presence.GameKey.TitleID = types.NewPrimitiveU64(0)
		//bella.Presence.GameKey.TitleVersion = 64
		bella.Presence.GameKey.TitleVersion = types.NewPrimitiveU16(0)

		bella.Status.Unknown = types.NewPrimitiveU8(0)
		bella.Status.Contents = types.NewString("test")
		bella.Status.LastChanged = types.NewDateTime(0)

		friendList.Append(bella)
	}

	rmcResponseStream := nex.NewByteStreamOut(globals.SecureEndpoint.LibraryVersions(), globals.SecureEndpoint.ByteStreamSettings())

	principalPreference.WriteTo(rmcResponseStream)
	comment.WriteTo(rmcResponseStream)
	friendList.WriteTo(rmcResponseStream)
	friendRequestsOut.WriteTo(rmcResponseStream)
	friendRequestsIn.WriteTo(rmcResponseStream)
	blockList.WriteTo(rmcResponseStream)
	types.NewPrimitiveBool(false).WriteTo(rmcResponseStream) // * Unknown
	notifications.WriteTo(rmcResponseStream)
	types.NewPrimitiveBool(false).WriteTo(rmcResponseStream) // * Unknown

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(globals.SecureEndpoint, rmcResponseBody)
	rmcResponse.ProtocolID = friends_wiiu.ProtocolID
	rmcResponse.MethodID = friends_wiiu.MethodUpdateAndGetAllInformation
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
