package utility

import (
	"context"
	"fmt"
	"time"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	account_management_types "github.com/PretendoNetwork/nex-protocols-go/v2/account-management/types"
	pb "github.com/PretendoNetwork/grpc/go/account/v2"
	"google.golang.org/grpc/metadata"

	"github.com/PretendoNetwork/friends/globals"
)

// ValidateNintendoCreateAccountToken validates the given Pretendo Network token for account creation
func ValidateNintendoCreateAccountToken(token types.DataHolder) (types.PID, *nex.Error) {
	var tokenBase64 string

	tokenDataType := token.Object.DataObjectID().(types.String)

	switch tokenDataType {
	case "NintendoCreateAccountData": // * Wii U
		nintendoCreateAccountData := token.Object.Copy().(account_management_types.NintendoCreateAccountData)

		tokenBase64 = string(nintendoCreateAccountData.Token)
	case "AccountExtraInfo": // * 3DS
		accountExtraInfo := token.Object.Copy().(account_management_types.AccountExtraInfo)

		tokenBase64 = string(accountExtraInfo.NEXToken)
	default:
		globals.Logger.Errorf("Invalid token data type %s!", tokenDataType)
		return 0, nex.NewError(nex.ResultCodes.Authentication.ValidationFailed, fmt.Sprintf("Invalid token data type %s!", tokenDataType))
	}

	ctx := metadata.NewOutgoingContext(context.Background(), common_globals.GRPCAccountCommonMetadata)

	response, err := common_globals.GRPCAccountClient.ExchangeNEXTokenForUserData(ctx, &pb.ExchangeNEXTokenForUserDataRequest{
		GameServerIds: []string{"00003200"},
		Token:         tokenBase64,
	})
	if err != nil {
		return 0, nex.NewError(nex.ResultCodes.Authentication.ValidationFailed, err.Error())
	}

	// * The account server database separates all the token types into their own
	// * collections, so a non-NEX token (even if valid) should still return no
	// * data here. But sanity the types check anyway just in case
	if response.TokenInfo.TokenType != 3 { // * 3 = NEX
		return 0, nex.NewError(nex.ResultCodes.Authentication.ValidationFailed, "Invalid token")
	}

	if response.TokenInfo.SystemType != 1 && response.TokenInfo.SystemType != 2 { // * 1 = WUP, 2 = CTR
		return 0, nex.NewError(nex.ResultCodes.Authentication.ValidationFailed, "Invalid token")
	}

	// * If the token is expired, the account server database will have deleted it,
	// * but sanity check anyway just in case
	if response.TokenInfo.ExpireTime != nil && response.TokenInfo.ExpireTime.AsTime().Before(time.Now()) {
		return 0, nex.NewError(nex.ResultCodes.Authentication.TokenExpired, "Token expired")
	}

	// PID isn't checked since account creation is done with a guest account

	if response.NexAccount.AccessLevel < 0 {
		return 0, nex.NewError(nex.ResultCodes.RendezVous.AccountDisabled, fmt.Sprintf("Account %d is banned", response.NexAccount.Pid))
	}

	return types.NewPID(uint64(response.NexAccount.Pid)), nil
}
