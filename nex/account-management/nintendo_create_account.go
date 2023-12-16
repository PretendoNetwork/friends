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
	nex "github.com/PretendoNetwork/nex-go"
	account_management "github.com/PretendoNetwork/nex-protocols-go/account-management"
	account_management_types "github.com/PretendoNetwork/nex-protocols-go/account-management/types"
)

func NintendoCreateAccount(err error, packet nex.PacketInterface, callID uint32, strPrincipalName string, strKey string, uiGroups uint32, strEmail string, oAuthData *nex.DataHolder) (*nex.RMCMessage, uint32) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.Errors.Core.InvalidArgument
	}

	var tokenBase64 string

	oAuthDataType := oAuthData.TypeName()

	switch oAuthDataType {
	case "NintendoCreateAccountData": // * Wii U
		nintendoCreateAccountData := oAuthData.ObjectData().(*account_management_types.NintendoCreateAccountData)

		tokenBase64 = nintendoCreateAccountData.Token
	case "AccountExtraInfo": // * 3DS
		accountExtraInfo := oAuthData.ObjectData().(*account_management_types.AccountExtraInfo)

		tokenBase64 = accountExtraInfo.NEXToken
		tokenBase64 = strings.Replace(tokenBase64, ".", "+", -1)
		tokenBase64 = strings.Replace(tokenBase64, "-", "/", -1)
		tokenBase64 = strings.Replace(tokenBase64, "*", "=", -1)
	default:
		globals.Logger.Errorf("Invalid oAuthData data type %s!", oAuthDataType)
		return nil, nex.Errors.Authentication.TokenParseError
	}

	encryptedToken, err := base64.StdEncoding.DecodeString(tokenBase64)
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.Errors.Authentication.TokenParseError
	}

	decryptedToken, err := utility.DecryptToken(encryptedToken)
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.Errors.Authentication.TokenParseError
	}

	pid := decryptedToken.UserPID

	pidByteArray := make([]byte, 4)
	binary.LittleEndian.PutUint32(pidByteArray, pid)

	mac := hmac.New(md5.New, []byte(strKey))
	_, err = mac.Write(pidByteArray)
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.Errors.Authentication.Unknown
	}

	pidHmac := hex.EncodeToString(mac.Sum(nil))

	rmcResponseStream := nex.NewStreamOut(globals.SecureServer)

	rmcResponseStream.WriteUInt32LE(pid)
	rmcResponseStream.WriteString(pidHmac)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(globals.SecureServer, rmcResponseBody)
	rmcResponse.ProtocolID = account_management.ProtocolID
	rmcResponse.MethodID = account_management.MethodNintendoCreateAccount
	rmcResponse.CallID = callID

	return rmcResponse, 0
}
