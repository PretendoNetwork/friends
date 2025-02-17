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

func (s *gRPCFriendsServer) GetUserFriendRequestsIncoming(ctx context.Context, in *pb.GetUserFriendRequestsIncomingRequest) (*pb.GetUserFriendRequestsIncomingResponse, error) {
	friendRequestsIn, err := database_wiiu.GetUserFriendRequestsIn(in.GetPid())
	if err != nil && err != database.ErrEmptyList {
		globals.Logger.Critical(err.Error())
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	friendRequests := make([]*pb.FriendRequest, 0, len(friendRequestsIn))

	if friendRequestsIn != nil {
		for _, friendRequestIn := range friendRequestsIn {
			friendRequest := &pb.FriendRequest{
				Id:        uint64(friendRequestIn.Message.FriendRequestID),
				Sender:    uint32(friendRequestIn.PrincipalInfo.PID),
				Recipient: in.GetPid(),
				Sent:      uint64(friendRequestIn.SentOn.Standard().Unix()),
				Expires:   uint64(friendRequestIn.Message.ExpiresOn.Standard().Unix()),
				Message:   string(friendRequestIn.Message.Message),
			}

			friendRequests = append(friendRequests, friendRequest)
		}
	}

	return &pb.GetUserFriendRequestsIncomingResponse{
		FriendRequests: friendRequests,
	}, nil
}
