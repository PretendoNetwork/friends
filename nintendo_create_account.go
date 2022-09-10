package main

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"strings"

	"github.com/PretendoNetwork/friends-secure/globals"
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func nintendoCreateAccount(err error, client *nex.Client, callID uint32, strPrincipalName string, strKey string, uiGroups uint32, strEmail string, oAuthData *nex.DataHolder) {
	if err != nil {
		// TODO: Handle error
		globals.Logger.Critical(err.Error())
	}

	var tokenBase64 string

	oAuthDataType := oAuthData.TypeName()

	if oAuthDataType == "NintendoCreateAccountData" { // Wii U
		nintendoCreateAccountData := oAuthData.ObjectData().(*nexproto.NintendoCreateAccountData)

		tokenBase64 = nintendoCreateAccountData.Token
	} else if oAuthDataType == "AccountExtraInfo" { // 3DS
		accountExtraInfo := oAuthData.ObjectData().(*nexproto.AccountExtraInfo)

		tokenBase64 = accountExtraInfo.NEXToken
		tokenBase64 = strings.Replace(tokenBase64, ".", "+", -1)
		tokenBase64 = strings.Replace(tokenBase64, "-", "/", -1)
		tokenBase64 = strings.Replace(tokenBase64, "*", "=", -1)
	}

	encryptedToken, _ := base64.StdEncoding.DecodeString(tokenBase64)

	decryptedToken, err := decryptToken(encryptedToken)
	if err != nil {
		// TODO: Handle error
		globals.Logger.Critical(err.Error())
	}

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

	rmcResponse := nex.NewRMCResponse(nexproto.AccountManagementProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.AccountManagementMethodNintendoCreateAccount, rmcResponseBody)

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
