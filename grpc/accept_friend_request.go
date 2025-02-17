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

func (s *gRPCFriendsServer) AcceptFriendRequest(ctx context.Context, in *pb.AcceptFriendRequestRequest) (*pb.AcceptFriendRequestResponse, error) {

	_, err := database_wiiu.AcceptFriendRequestAndReturnFriendInfo(in.GetFriendRequestId())
	if err != nil {
		if err == database.ErrFriendRequestNotFound {
			return &pb.AcceptFriendRequestResponse{
				Success: false,
			}, status.Errorf(codes.NotFound, "friend request not found")
		}

		if err == database.ErrPIDNotFound {
			return &pb.AcceptFriendRequestResponse{
				Success: false,
			}, status.Errorf(codes.FailedPrecondition, "friend request has invalid PID")
		}

		globals.Logger.Critical(err.Error())
		return &pb.AcceptFriendRequestResponse{
			Success: false,
		}, status.Errorf(codes.Internal, "internal server error")
	}

	return &pb.AcceptFriendRequestResponse{
		Success: err == nil,
	}, nil
}
