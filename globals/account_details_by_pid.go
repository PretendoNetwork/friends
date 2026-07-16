package globals

import (
	"context"
	"strconv"

	pb "github.com/PretendoNetwork/grpc/go/account/v2"
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	"google.golang.org/grpc/metadata"
)

func AccountDetailsByPID(pid types.PID) (*nex.Account, *nex.Error) {
	if pid.Equals(AuthenticationServerAccount.PID) {
		return AuthenticationServerAccount, nil
	}

	if pid.Equals(SecureServerAccount.PID) {
		return SecureServerAccount, nil
	}

	if pid.Equals(GuestAccount.PID) {
		return GuestAccount, nil
	}

	ctx := metadata.NewOutgoingContext(context.Background(), common_globals.GRPCAccountCommonMetadata)

	response, err := common_globals.GRPCAccountClient.GetNEXPassword(ctx, &pb.GetNEXPasswordRequest{Pid: uint32(pid)})
	if err != nil {
		Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.RendezVous.InvalidPID, "Invalid PID")
	}

	username := strconv.Itoa(int(pid))
	account := nex.NewAccount(pid, username, response.Password, false)

	return account, nil
}
