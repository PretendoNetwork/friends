package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	database_wiiu "github.com/PretendoNetwork/friends/database/wiiu"
	"github.com/PretendoNetwork/friends/globals"
	pb "github.com/PretendoNetwork/grpc/go/friends"
	"github.com/PretendoNetwork/nex-go/v2/types"
)

func (s *gRPCFriendsServer) SendUserFriendRequest(ctx context.Context, in *pb.SendUserFriendRequestRequest) (*pb.SendUserFriendRequestResponse, error) {
	sender := in.GetSender()
	recipient := in.GetRecipient()

	// * Checking that the recipient allows friend requests
	recipientPrincipalPreferences, err := database_wiiu.GetUserPrincipalPreference(recipient)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return &pb.SendUserFriendRequestResponse{
			Success: false,
		}, status.Errorf(codes.Internal, "internal server error")
	}

	if recipientPrincipalPreferences.BlockFriendRequests == true {
		// * Do not allow a user with friend requests off to receive a friend request
		return &pb.SendUserFriendRequestResponse{
			Success: false,
		}, status.Errorf(codes.Unavailable, "Provisional friend requests through gRPC are unavailable.")
	}

	print("check is fucked")
	currentTimestamp := time.Now()
	expireTimestamp := currentTimestamp.Add(time.Hour * 24 * 29)

	sentTime := types.NewDateTime(0)
	expireTime := types.NewDateTime(0)

	sentTime.FromTimestamp(currentTimestamp)
	expireTime.FromTimestamp(expireTimestamp)

	message := in.GetMessage()

	id, err := database_wiiu.SaveFriendRequest(sender, recipient, uint64(sentTime), uint64(expireTime), message)
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
