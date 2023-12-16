package globals

import (
	"context"

	pb "github.com/PretendoNetwork/grpc-go/account"
	"github.com/PretendoNetwork/nex-go"
	"google.golang.org/grpc/metadata"
)

func PasswordFromPID(pid *nex.PID) (string, uint32) {
	ctx := metadata.NewOutgoingContext(context.Background(), GRPCAccountCommonMetadata)

	response, err := GRPCAccountClient.GetNEXPassword(ctx, &pb.GetNEXPasswordRequest{Pid: pid.LegacyValue()})
	if err != nil {
		Logger.Error(err.Error())
		return "", nex.Errors.RendezVous.InvalidUsername
	}

	return response.Password, 0
}
