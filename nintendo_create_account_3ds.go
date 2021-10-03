package main

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"strings"

	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func nintendoCreateAccount3DS(err error, client *nex.Client, callID uint32, strPrincipalName string, strKey string, uiGroups uint32, strEmail string, oAuthData *nexproto.AccountExtraInfo) {
	if err != nil {
		panic(err)
	}

	tokenBase64 := oAuthData.NEXToken
	tokenBase64 = strings.Replace(tokenBase64, ".", "+", -1)
	tokenBase64 = strings.Replace(tokenBase64, "-", "/", -1)
	tokenBase64 = strings.Replace(tokenBase64, "*", "=", -1)

	encryptedToken, _ := base64.StdEncoding.DecodeString(tokenBase64)

	decryptedToken, err := decryptToken(encryptedToken)
	if err != nil {
		panic(err)
	}

	pid := decryptedToken.UserPID

	pidByteArray := make([]byte, 4)
	binary.LittleEndian.PutUint32(pidByteArray, pid)

	mac := hmac.New(md5.New, []byte(strKey))
	mac.Write(pidByteArray)

	pidHmac := hex.EncodeToString(mac.Sum(nil))

	rmcResponseStream := nex.NewStreamOut(nexServer)

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

	nexServer.Send(responsePacket)
}
