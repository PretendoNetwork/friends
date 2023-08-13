package globals

import (
	"context"

	pb "github.com/PretendoNetwork/grpc-go/account"
	"github.com/PretendoNetwork/nex-protocols-go/globals"
	"google.golang.org/grpc/metadata"
)

func GetUserData(pid uint32) (*pb.GetUserDataResponse, error) {
	ctx := metadata.NewOutgoingContext(context.Background(), GRPCAccountCommonMetadata)

	response, err := GRPCAccountClient.GetUserData(ctx, &pb.GetUserDataRequest{Pid: pid})
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, err
	}

	return response, nil
}
