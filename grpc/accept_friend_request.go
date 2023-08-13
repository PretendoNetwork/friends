package grpc

import (
	"context"
	"database/sql"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	database_wiiu "github.com/PretendoNetwork/friends/database/wiiu"
	pb "github.com/PretendoNetwork/grpc-go/friends"
)

func (s *gRPCFriendsServer) AcceptFriendRequest(ctx context.Context, in *pb.AcceptFriendRequestRequest) (*pb.AcceptFriendRequestResponse, error) {

	friendInfo, err := database_wiiu.AcceptFriendRequestAndReturnFriendInfo(in.GetFriendRequestId())
	if err == sql.ErrNoRows {
		return &pb.AcceptFriendRequestResponse{
			Success: false,
		}, status.Errorf(codes.NotFound, "friend request not found")
	}

	return &pb.AcceptFriendRequestResponse{
		Success: friendInfo != nil,
	}, nil
}
