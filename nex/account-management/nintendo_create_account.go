package nex_account_management

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"strings"

	"github.com/PretendoNetwork/friends/globals"
	"github.com/PretendoNetwork/friends/utility"
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	account_management "github.com/PretendoNetwork/nex-protocols-go/v2/account-management"
	account_management_types "github.com/PretendoNetwork/nex-protocols-go/v2/account-management/types"
)

func NintendoCreateAccount(err error, packet nex.PacketInterface, callID uint32, strPrincipalName *types.String, strKey *types.String, uiGroups *types.PrimitiveU32, strEmail *types.String, oAuthData *types.AnyDataHolder) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "") // TODO - Add error message
	}

	var tokenBase64 string

	oAuthDataType := oAuthData.TypeName.Value

	switch oAuthDataType {
	case "NintendoCreateAccountData": // * Wii U
		nintendoCreateAccountData := oAuthData.ObjectData.Copy().(*account_management_types.NintendoCreateAccountData)

		tokenBase64 = nintendoCreateAccountData.Token.Value
	case "AccountExtraInfo": // * 3DS
		accountExtraInfo := oAuthData.ObjectData.Copy().(*account_management_types.AccountExtraInfo)

		tokenBase64 = accountExtraInfo.NEXToken.Value
		tokenBase64 = strings.Replace(tokenBase64, ".", "+", -1)
		tokenBase64 = strings.Replace(tokenBase64, "-", "/", -1)
		tokenBase64 = strings.Replace(tokenBase64, "*", "=", -1)
	default:
		globals.Logger.Errorf("Invalid oAuthData data type %s!", oAuthDataType)
		return nil, nex.NewError(nex.ResultCodes.Authentication.TokenParseError, "") // TODO - Add error message
	}

	encryptedToken, err := base64.StdEncoding.DecodeString(tokenBase64)
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Authentication.TokenParseError, "") // TODO - Add error message
	}

	decryptedToken, err := utility.DecryptToken(encryptedToken)
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Authentication.TokenParseError, "") // TODO - Add error message
	}

	pid := types.NewPID(uint64(decryptedToken.UserPID))

	pidByteArray := make([]byte, 4)
	binary.LittleEndian.PutUint32(pidByteArray, pid.LegacyValue())

	mac := hmac.New(md5.New, []byte(strKey.Value))
	_, err = mac.Write(pidByteArray)
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Authentication.Unknown, "") // TODO - Add error message
	}

	pidHmac := types.NewString(hex.EncodeToString(mac.Sum(nil)))

	rmcResponseStream := nex.NewByteStreamOut(globals.SecureEndpoint.LibraryVersions(), globals.SecureEndpoint.ByteStreamSettings())

	pid.WriteTo(rmcResponseStream)
	pidHmac.WriteTo(rmcResponseStream)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(globals.SecureEndpoint, rmcResponseBody)
	rmcResponse.ProtocolID = account_management.ProtocolID
	rmcResponse.MethodID = account_management.MethodNintendoCreateAccount
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
