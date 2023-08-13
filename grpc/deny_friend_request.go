package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/PretendoNetwork/friends/database"
	database_wiiu "github.com/PretendoNetwork/friends/database/wiiu"
	"github.com/PretendoNetwork/friends/globals"
	pb "github.com/PretendoNetwork/grpc-go/friends"
)

func (s *gRPCFriendsServer) DenyFriendRequest(ctx context.Context, in *pb.DenyFriendRequestRequest) (*pb.DenyFriendRequestResponse, error) {
	err := database_wiiu.SetFriendRequestDenied(in.GetFriendRequestId())
	if err != nil {
		if err == database.ErrFriendRequestNotFound {
			return &pb.DenyFriendRequestResponse{
				Success: false,
			}, status.Errorf(codes.NotFound, "friend request not found")
		}

		globals.Logger.Critical(err.Error())
		return &pb.DenyFriendRequestResponse{
			Success: false,
		}, status.Errorf(codes.Internal, "internal server error")
	}

	return &pb.DenyFriendRequestResponse{
		Success: true,
	}, nil
}
