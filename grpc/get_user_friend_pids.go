package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/PretendoNetwork/friends/database"
	database_3ds "github.com/PretendoNetwork/friends/database/3ds"
	database_wiiu "github.com/PretendoNetwork/friends/database/wiiu"
	"github.com/PretendoNetwork/friends/globals"
	pb "github.com/PretendoNetwork/grpc-go/friends"
	friends_3ds_types "github.com/PretendoNetwork/nex-protocols-go/v2/friends-3ds/types"
)

func (s *gRPCFriendsServer) GetUserFriendPIDs(ctx context.Context, in *pb.GetUserFriendPIDsRequest) (*pb.GetUserFriendPIDsResponse, error) {
	var pids []uint32
	var err error

	// * Try Wii U database first
	pids, err = database_wiiu.GetUserFriendPIDs(in.GetPid())
	if err != nil && err != database.ErrEmptyList {
		globals.Logger.Critical(err.Error())
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	if len(pids) > 0 {
		return &pb.GetUserFriendPIDsResponse{
			Pids: pids,
		}, nil
	}

	// * If no PIDs are given, try with 3DS database instead
	relationships, err := database_3ds.GetUserFriends(in.GetPid())
	if err != nil && err != database.ErrEmptyList {
		globals.Logger.Critical(err.Error())
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	if relationships != nil {
		relationships.Each(func(i int, relationship *friends_3ds_types.FriendRelationship) bool {
			// * Only add complete relationships to the list
			if relationship.RelationshipType.Value == 1 {
				pids = append(pids, relationship.PID.LegacyValue())
			}

			return false
		})
	}

	return &pb.GetUserFriendPIDsResponse{
		Pids: pids,
	}, nil
}
