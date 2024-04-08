package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/PretendoNetwork/friends/database"
	database_wiiu "github.com/PretendoNetwork/friends/database/wiiu"
	"github.com/PretendoNetwork/friends/globals"
	pb "github.com/PretendoNetwork/grpc-go/friends"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/v2/friends-wiiu/types"
)

func (s *gRPCFriendsServer) GetUserFriendRequestsIncoming(ctx context.Context, in *pb.GetUserFriendRequestsIncomingRequest) (*pb.GetUserFriendRequestsIncomingResponse, error) {
	friendRequestsIn, err := database_wiiu.GetUserFriendRequestsIn(in.GetPid())
	if err != nil && err != database.ErrEmptyList {
		globals.Logger.Critical(err.Error())
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	friendRequests := make([]*pb.FriendRequest, 0, friendRequestsIn.Length())

	if friendRequestsIn != nil {
		friendRequestsIn.Each(func(i int, friendRequestIn *friends_wiiu_types.FriendRequest) bool {
			friendRequest := &pb.FriendRequest{
				Id:        friendRequestIn.Message.FriendRequestID.Value,
				Sender:    friendRequestIn.PrincipalInfo.PID.LegacyValue(),
				Recipient: in.GetPid(),
				Sent:      uint64(friendRequestIn.SentOn.Standard().Unix()),
				Expires:   uint64(friendRequestIn.Message.ExpiresOn.Standard().Unix()),
				Message:   friendRequestIn.Message.Message.Value,
			}

			friendRequests = append(friendRequests, friendRequest)

			return false
		})
	}

	return &pb.GetUserFriendRequestsIncomingResponse{
		FriendRequests: friendRequests,
	}, nil
}
