package globals

import (
	"context"

	pb "github.com/PretendoNetwork/grpc-go/account"
	"google.golang.org/grpc/metadata"
)

func GetUserData(pid uint32) (*pb.GetUserDataResponse, error) {
	ctx := metadata.NewOutgoingContext(context.Background(), GRPCAccountCommonMetadata)

	response, err := GRPCAccountClient.GetUserData(ctx, &pb.GetUserDataRequest{Pid: pid})
	if err != nil {
		return nil, err
	}

	return response, nil
}
