package grpc

import (
	"context"

	database_wiiu "github.com/PretendoNetwork/friends-secure/database/wiiu"
	pb "github.com/PretendoNetwork/grpc-go/friends"
)

func (s *gRPCFriendsServer) AcceptFriendRequest(ctx context.Context, in *pb.AcceptFriendRequestRequest) (*pb.AcceptFriendRequestResponse, error) {

	friendInfo := database_wiiu.AcceptFriendRequestAndReturnFriendInfo(in.GetFriendRequestId())

	return &pb.AcceptFriendRequestResponse{
		Success: friendInfo != nil,
	}, nil
}
