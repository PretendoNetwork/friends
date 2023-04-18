package grpc

import (
	"context"

	database_wiiu "github.com/PretendoNetwork/friends-secure/database/wiiu"
	pb "github.com/PretendoNetwork/grpc-go/friends"
)

func (s *gRPCFriendsServer) DenyFriendRequest(ctx context.Context, in *pb.DenyFriendRequestRequest) (*pb.DenyFriendRequestResponse, error) {
	// TODO - Make this return an error
	database_wiiu.SetFriendRequestDenied(in.GetFriendRequestId())

	return &pb.DenyFriendRequestResponse{
		Success: true,
	}, nil
}
