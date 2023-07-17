package nex_account_management

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"strings"

	"github.com/PretendoNetwork/friends-secure/globals"
	"github.com/PretendoNetwork/friends-secure/utility"
	nex "github.com/PretendoNetwork/nex-go"
	account_management "github.com/PretendoNetwork/nex-protocols-go/account-management"
)

func NintendoCreateAccount(err error, client *nex.Client, callID uint32, strPrincipalName string, strKey string, uiGroups uint32, strEmail string, oAuthData *nex.DataHolder) {
	if err != nil {
		// TODO: Handle error
		globals.Logger.Critical(err.Error())
	}

	rmcResponse := nex.NewRMCResponse(account_management.ProtocolID, callID)

	var tokenBase64 string

	oAuthDataType := oAuthData.TypeName()

	if oAuthDataType == "NintendoCreateAccountData" { // Wii U
		nintendoCreateAccountData := oAuthData.ObjectData().(*account_management.NintendoCreateAccountData)

		tokenBase64 = nintendoCreateAccountData.Token
	} else if oAuthDataType == "AccountExtraInfo" { // 3DS
		accountExtraInfo := oAuthData.ObjectData().(*account_management.AccountExtraInfo)

		tokenBase64 = accountExtraInfo.NEXToken
		tokenBase64 = strings.Replace(tokenBase64, ".", "+", -1)
		tokenBase64 = strings.Replace(tokenBase64, "-", "/", -1)
		tokenBase64 = strings.Replace(tokenBase64, "*", "=", -1)
	}

	encryptedToken, _ := base64.StdEncoding.DecodeString(tokenBase64)

	decryptedToken, err := utility.DecryptToken(encryptedToken)
	if err != nil {
		globals.Logger.Error(err.Error())
		rmcResponse.SetError(nex.Errors.Authentication.TokenParseError)
	} else {
		pid := decryptedToken.UserPID

		pidByteArray := make([]byte, 4)
		binary.LittleEndian.PutUint32(pidByteArray, pid)

		mac := hmac.New(md5.New, []byte(strKey))
		mac.Write(pidByteArray)

		pidHmac := hex.EncodeToString(mac.Sum(nil))

		rmcResponseStream := nex.NewStreamOut(globals.NEXServer)

		rmcResponseStream.WriteUInt32LE(pid)
		rmcResponseStream.WriteString(pidHmac)

		rmcResponseBody := rmcResponseStream.Bytes()

		rmcResponse.SetSuccess(account_management.MethodNintendoCreateAccount, rmcResponseBody)
	}

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewPacketV0(client, nil)

	responsePacket.SetVersion(0)
	responsePacket.SetSource(0xA1)
	responsePacket.SetDestination(0xAF)
	responsePacket.SetType(nex.DataPacket)
	responsePacket.SetPayload(rmcResponseBytes)

	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)

	globals.NEXServer.Send(responsePacket)
}
