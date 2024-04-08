package globals

import (
	"context"
	"fmt"
	"strconv"

	pb "github.com/PretendoNetwork/grpc-go/account"
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"google.golang.org/grpc/metadata"
)

func AccountDetailsByUsername(username string) (*nex.Account, *nex.Error) {
	if username == AuthenticationEndpoint.ServerAccount.Username {
		return AuthenticationEndpoint.ServerAccount, nil
	}

	if username == SecureEndpoint.ServerAccount.Username {
		return SecureEndpoint.ServerAccount, nil
	}

	if username == GuestAccount.Username {
		return GuestAccount, nil
	}

	// TODO - This is fine for our needs, but not for servers which use non-PID usernames?
	pid, err := strconv.Atoi(username)
	if err != nil {
		fmt.Println(1)
		fmt.Println(err)
		return nil, nex.NewError(nex.ResultCodes.RendezVous.InvalidUsername, "Invalid username")
	}

	// * Trying to use AccountDetailsByPID here led to weird nil checks?
	// * Would always return an error even when it shouldn't.
	// TODO - Look into this more

	ctx := metadata.NewOutgoingContext(context.Background(), GRPCAccountCommonMetadata)

	response, err := GRPCAccountClient.GetNEXPassword(ctx, &pb.GetNEXPasswordRequest{Pid: uint32(pid)})
	if err != nil {
		Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.RendezVous.InvalidPID, "Invalid PID")
	}

	account := nex.NewAccount(types.NewPID(uint64(pid)), username, response.Password)

	return account, nil
}
