package friends_wiiu

import (
	"time"

	database_wiiu "github.com/PretendoNetwork/friends-secure/database/wiiu"
	"github.com/PretendoNetwork/friends-secure/globals"
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func DenyFriendRequest(err error, client *nex.Client, callID uint32, id uint64) {
	database_wiiu.SetFriendRequestDenied(id)

	senderPID, _ := database_wiiu.GetPIDsByFriendRequestID(id)
	database_wiiu.SetUserBlocked(client.PID(), senderPID, 0, 0)

	info := database_wiiu.GetUserInfoByPID(senderPID)

	date := nex.NewDateTime(0)
	date.FromTimestamp(time.Now())

	// Create a new blacklist principal for the client, as unlike AddBlacklist they don't send one to us.
	blacklistPrincipal := nexproto.NewBlacklistedPrincipal()

	blacklistPrincipal.PrincipalBasicInfo = info
	blacklistPrincipal.GameKey = nexproto.NewGameKey()
	blacklistPrincipal.BlackListedSince = date

	rmcResponseStream := nex.NewStreamOut(globals.NEXServer)

	rmcResponseStream.WriteStructure(blacklistPrincipal)

	rmcResponseBody := rmcResponseStream.Bytes()

	// Build response packet
	rmcResponse := nex.NewRMCResponse(nexproto.FriendsWiiUProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.FriendsWiiUMethodDenyFriendRequest, rmcResponseBody)

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
