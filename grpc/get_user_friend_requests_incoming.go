package grpc

import (
	"context"

	database_wiiu "github.com/PretendoNetwork/friends-secure/database/wiiu"
	pb "github.com/PretendoNetwork/grpc-go/friends"
)

func (s *gRPCFriendsServer) GetUserFriendRequestsIncoming(ctx context.Context, in *pb.GetUserFriendRequestsIncomingRequest) (*pb.GetUserFriendRequestsIncomingResponse, error) {

	friendRequestsIn := database_wiiu.GetUserFriendRequestsIn(in.GetPid())

	friendRequests := make([]*pb.FriendRequest, 0, len(friendRequestsIn))

	for i := 0; i < len(friendRequestsIn); i++ {
		friendRequest := &pb.FriendRequest{
			Id:        friendRequestsIn[i].Message.FriendRequestID,
			Sender:    friendRequestsIn[i].PrincipalInfo.PID,
			Recipient: in.GetPid(),
			Sent:      uint64(friendRequestsIn[i].SentOn.Standard().Unix()),
			Expires:   uint64(friendRequestsIn[i].Message.ExpiresOn.Standard().Unix()),
			Message:   friendRequestsIn[i].Message.Message,
		}

		friendRequests = append(friendRequests, friendRequest)
	}

	return &pb.GetUserFriendRequestsIncomingResponse{
		FriendRequests: friendRequests,
	}, nil
}
