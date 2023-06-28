package grpc

import (
	"context"

	database_wiiu "github.com/PretendoNetwork/friends-secure/database/wiiu"
	database_3ds "github.com/PretendoNetwork/friends-secure/database/3ds"
	pb "github.com/PretendoNetwork/grpc-go/friends"
)

func (s *gRPCFriendsServer) GetUserFriendPIDs(ctx context.Context, in *pb.GetUserFriendPIDsRequest) (*pb.GetUserFriendPIDsResponse, error) {
	var pids []uint32

	// * Try Wii U database first
	pids = database_wiiu.GetUserFriendPIDs(in.GetPid())

	if len(pids) > 0 {
		return &pb.GetUserFriendPIDsResponse{
			Pids: pids,
		}, nil
	}

	// * If no PIDs are given, try with 3DS database instead
	relationships := database_3ds.GetUserFriends(in.GetPid())

	for _, relationship := range relationships {
		// * Only add complete relationships to the list
		if relationship.RelationshipType == 1 {
			pids = append(pids, relationship.PID)
		}
	}

	return &pb.GetUserFriendPIDsResponse{
		Pids: pids,
	}, nil
}
