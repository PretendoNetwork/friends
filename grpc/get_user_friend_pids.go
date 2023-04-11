package grpc

import (
	"context"

	database_wiiu "github.com/PretendoNetwork/friends-secure/database/wiiu"
	pb "github.com/PretendoNetwork/grpc-go/friends"
)

func (s *gRPCFriendsServer) GetUserFriendPIDs(ctx context.Context, in *pb.GetUserFriendPIDsRequest) (*pb.GetUserFriendPIDsResponse, error) {

	pids := database_wiiu.GetUserFriendPIDs(in.GetPid())

	return &pb.GetUserFriendPIDsResponse{
		Pids: pids,
	}, nil
}
