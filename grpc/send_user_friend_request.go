package grpc

import (
	"context"
	"time"

	database_wiiu "github.com/PretendoNetwork/friends-secure/database/wiiu"
	pb "github.com/PretendoNetwork/grpc-go/friends"
	nex "github.com/PretendoNetwork/nex-go"
)

func (s *gRPCFriendsServer) SendUserFriendRequest(ctx context.Context, in *pb.SendUserFriendRequestRequest) (*pb.SendUserFriendRequestResponse, error) {
	sender := in.GetSender()
	recipient := in.GetRecipient()

	currentTimestamp := time.Now()
	expireTimestamp := currentTimestamp.Add(time.Hour * 24 * 29)

	sentTime := nex.NewDateTime(0)
	expireTime := nex.NewDateTime(0)

	sentTime.FromTimestamp(currentTimestamp)
	expireTime.FromTimestamp(expireTimestamp)

	message := in.GetMessage()

	id := database_wiiu.SaveFriendRequest(sender, recipient, sentTime.Value(), expireTime.Value(), message)

	return &pb.SendUserFriendRequestResponse{
		Success: id == 0,
	}, nil
}
