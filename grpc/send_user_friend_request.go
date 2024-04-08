package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	database_wiiu "github.com/PretendoNetwork/friends/database/wiiu"
	"github.com/PretendoNetwork/friends/globals"
	pb "github.com/PretendoNetwork/grpc-go/friends"
	"github.com/PretendoNetwork/nex-go/v2/types"
)

func (s *gRPCFriendsServer) SendUserFriendRequest(ctx context.Context, in *pb.SendUserFriendRequestRequest) (*pb.SendUserFriendRequestResponse, error) {
	sender := in.GetSender()
	recipient := in.GetRecipient()

	currentTimestamp := time.Now()
	expireTimestamp := currentTimestamp.Add(time.Hour * 24 * 29)

	sentTime := types.NewDateTime(0)
	expireTime := types.NewDateTime(0)

	sentTime.FromTimestamp(currentTimestamp)
	expireTime.FromTimestamp(expireTimestamp)

	message := in.GetMessage()

	id, err := database_wiiu.SaveFriendRequest(sender, recipient, sentTime.Value(), expireTime.Value(), message)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return &pb.SendUserFriendRequestResponse{
			Success: false,
		}, status.Errorf(codes.Internal, "internal server error")
	}

	return &pb.SendUserFriendRequestResponse{
		Success: id != 0,
	}, nil
}
