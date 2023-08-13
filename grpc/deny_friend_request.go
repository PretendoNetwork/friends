package grpc

import (
	"context"

	database_wiiu "github.com/PretendoNetwork/friends/database/wiiu"
	pb "github.com/PretendoNetwork/grpc-go/friends"
)

func (s *gRPCFriendsServer) DenyFriendRequest(ctx context.Context, in *pb.DenyFriendRequestRequest) (*pb.DenyFriendRequestResponse, error) {
	err := database_wiiu.SetFriendRequestDenied(in.GetFriendRequestId())

	return &pb.DenyFriendRequestResponse{
		Success: err == nil,
	}, nil
}
