package friends_wiiu

import (
	"time"

	database_wiiu "github.com/PretendoNetwork/friends-secure/database/wiiu"
	"github.com/PretendoNetwork/friends-secure/globals"
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func AddBlacklist(err error, client *nex.Client, callID uint32, blacklistPrincipal *nexproto.BlacklistedPrincipal) {
	currentBlacklistPrincipal := blacklistPrincipal

	senderPID := currentBlacklistPrincipal.PrincipalBasicInfo.PID
	titleID := currentBlacklistPrincipal.GameKey.TitleID
	titleVersion := currentBlacklistPrincipal.GameKey.TitleVersion

	date := nex.NewDateTime(0)
	date.FromTimestamp(time.Now())

	userInfo := database_wiiu.GetUserInfoByPID(currentBlacklistPrincipal.PrincipalBasicInfo.PID)

	if userInfo == nil {
		rmcResponse := nex.NewRMCResponse(nexproto.FriendsWiiUProtocolID, callID)
		rmcResponse.SetError(nex.Errors.FPD.FriendNotExists) // TODO: Not sure if this is the correct error.

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

		return
	}

	currentBlacklistPrincipal.PrincipalBasicInfo = userInfo
	currentBlacklistPrincipal.BlackListedSince = date

	database_wiiu.SetUserBlocked(client.PID(), senderPID, titleID, titleVersion)

	rmcResponseStream := nex.NewStreamOut(globals.NEXServer)

	rmcResponseStream.WriteStructure(blacklistPrincipal)

	rmcResponseBody := rmcResponseStream.Bytes()

	// Build response packet
	rmcResponse := nex.NewRMCResponse(nexproto.FriendsWiiUProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.FriendsWiiUMethodAddBlackList, rmcResponseBody)

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