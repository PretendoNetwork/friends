package grpc

import (
	"context"
	"database/sql"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/PretendoNetwork/friends/database"
	database_3ds "github.com/PretendoNetwork/friends/database/3ds"
	"github.com/PretendoNetwork/friends/globals"
	"github.com/PretendoNetwork/friends/utility"
	pb "github.com/PretendoNetwork/grpc/go/friends/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	friends_3ds_constants "github.com/PretendoNetwork/nex-protocols-go/v2/friends-3ds/constants"
	friends_3ds_types "github.com/PretendoNetwork/nex-protocols-go/v2/friends-3ds/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *gRPCFriendsV2Server) GetUserFriendsData3DS(ctx context.Context, in *pb.GetUserFriendsData3DSRequest) (*pb.GetUserFriendsData3DSResponse, error) {
	var friends []*pb.FriendInfo3DS
	friendList, err := database_3ds.GetUserFriends(in.Pid)
	if err != nil && err != database.ErrEmptyList {
		globals.Logger.Critical(err.Error())
		return &pb.GetUserFriendsData3DSResponse{
			Friends: friends,
		}, status.Errorf(codes.Internal, "internal server error")
	}

	friendPIDs := make([]uint32, len(friendList))

	for _, friend := range friendList {
		friendPIDs = append(friendPIDs, uint32(friend.PID))
	}

	friendInfoList, err := database_3ds.GetFriendPersistentInfos(uint32(in.Pid), friendPIDs)
	if err != nil && err != sql.ErrNoRows {
		globals.Logger.Critical(err.Error())
		return &pb.GetUserFriendsData3DSResponse{
			Friends: friends,
		}, status.Errorf(codes.Internal, "internal server error")
	}

	miiList, err := database_3ds.GetFriendMiis(friendPIDs)
	if err != nil && err != sql.ErrNoRows {
		globals.Logger.Critical(err.Error())
		return &pb.GetUserFriendsData3DSResponse{
			Friends: friends,
		}, status.Errorf(codes.Internal, "internal server error")
	}

	if globals.Config.EnableBella {
		bella := friends_3ds_types.NewFriendPersistentInfo()

		bella.PID = types.NewPID(1743126339)
		bella.Region = types.NewUInt8(0)
		bella.Country = types.NewUInt8(0)
		bella.Area = types.NewUInt8(0)
		bella.Language = types.NewUInt8(0)
		bella.Platform = types.NewUInt8(0)
		bella.GameKey.TitleID = 0x0005000010176900
		bella.GameKey.TitleVersion = types.NewUInt16(0)
		bella.Message = "Howdy Hey!"
		bella.MessageUpdatedAt = types.NewDateTime(0)
		bella.MiiModifiedAt = types.NewDateTime(0)
		bella.LastOnline = types.NewDateTime(0)

		mii := friends_3ds_types.NewMii()
		mii.Name = types.NewString("Bandwidth")
		mii.ProfanityFlag = types.NewBool(false)
		mii.CharacterSet = friends_3ds_constants.MiiCharacterSetJUE
		mii.MiiData = types.NewBuffer([]byte{
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

		friendMii := friends_3ds_types.NewFriendMii()
		friendMii.PID = types.NewPID(uint64(bella.PID))
		friendMii.Mii = mii
		friendMii.ModifiedAt = types.NewDateTime(0)

		friendInfoList = append(friendInfoList, bella)
		miiList = append(miiList, friendMii)
	}

	for _, friend := range friendInfoList {
		gameKey := &pb.GameKey{
			TitleId:      uint64(friend.GameKey.TitleID),
			TitleVersion: uint32(friend.GameKey.TitleVersion),
		}

		miiIndex := -1
		for index, mii := range miiList {
			if mii.PID == friend.PID {
				miiIndex = index
				break
			}
		}

		if miiIndex == -1 {
			continue
		}

		miiData := miiList[miiIndex]
		mii := &pb.Mii{
			Name:             string(miiData.Mii.Name),
			ProfanityFlag:    bool(miiData.Mii.ProfanityFlag),
			CharacterSet:     uint32(miiData.Mii.CharacterSet),
			MiiDataEncrypted: miiData.Mii.MiiData,
		}

		mii_data, err := utility.DecryptMiiData(miiData.Mii.MiiData)
		if err == nil {
			mii.MiiData = mii_data
		}

		friendMii := &pb.FriendMii{
			Pid:        uint32(miiData.PID),
			Mii:        mii,
			ModifiedAt: timestamppb.New(time.Unix(int64(miiData.ModifiedAt.Second()), 0)),
		}

		presence := &pb.NintendoPresence{}
		connectedUser, ok := globals.ConnectedUsers.Get(uint32(friend.PID))

		if ok && connectedUser != nil {
			presence.ChangedFlags = uint32(connectedUser.Presence.ChangedFlags)
			presence.GameKey = &pb.GameKey{
				TitleId:      uint64(connectedUser.Presence.GameKey.TitleID),
				TitleVersion: uint32(connectedUser.Presence.GameKey.TitleVersion),
			}
			presence.Message = string(connectedUser.Presence.Message)
			presence.JoinAvailableFlag = uint32(connectedUser.Presence.JoinAvailableFlag)
			presence.MatchmakeType = uint32(connectedUser.Presence.MatchmakeType)
			presence.JoinGameId = uint32(connectedUser.Presence.JoinGameID)
			presence.JoinGameMode = uint32(connectedUser.Presence.JoinGameMode)
			presence.OwnerPid = uint32(connectedUser.Presence.OwnerPID)
			presence.JoinGroupId = uint32(connectedUser.Presence.JoinGroupID)
			presence.ApplicationArg = connectedUser.Presence.ApplicationArg
		}

		info := &pb.FriendInfo3DS{
			Pid:              uint32(friend.PID),
			Region:           uint32(friend.Region),
			Country:          uint32(friend.Country),
			Area:             uint32(friend.Area),
			Language:         uint32(friend.Language),
			Platform:         uint32(friend.Platform),
			Presence:         presence,
			GameKey:          gameKey,
			Message:          string(friend.Message),
			MessageUpdatedAt: timestamppb.New(friend.MessageUpdatedAt.Standard()),
			MiiModifiedAt:    timestamppb.New(friend.MiiModifiedAt.Standard()),
			LastOnline:       timestamppb.New(friend.LastOnline.Standard()),
			Mii:              friendMii,
		}
		friends = append(friends, info)
	}

	return &pb.GetUserFriendsData3DSResponse{
		Friends: friends,
	}, nil
}
