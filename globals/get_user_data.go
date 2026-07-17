package globals

import (
	"context"

	pb "github.com/PretendoNetwork/grpc/go/account/v2"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	"google.golang.org/grpc/metadata"
)

func GetUserData(pid uint32) (*pb.GetUserDataResponse, error) {
	ctx := metadata.NewOutgoingContext(context.Background(), common_globals.GRPCAccountCommonMetadata)

	response, err := common_globals.GRPCAccountClient.GetUserData(ctx, &pb.GetUserDataRequest{Pid: pid})
	if err != nil {
		return nil, err
	}

	return response, nil
}
