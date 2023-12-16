package nex_friends_wiiu

import (
	"time"

	"github.com/PretendoNetwork/friends/database"
	database_wiiu "github.com/PretendoNetwork/friends/database/wiiu"
	"github.com/PretendoNetwork/friends/globals"
	"github.com/PretendoNetwork/friends/utility"
	nex "github.com/PretendoNetwork/nex-go"
	friends_wiiu "github.com/PretendoNetwork/nex-protocols-go/friends-wiiu"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/friends-wiiu/types"
)

func DenyFriendRequest(err error, packet nex.PacketInterface, callID uint32, id uint64) (*nex.RMCMessage, uint32) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.Errors.FPD.InvalidArgument
	}

	client := packet.Sender().(*nex.PRUDPClient)

	err = database_wiiu.SetFriendRequestDenied(id)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return nil, nex.Errors.FPD.Unknown
	}

	senderPID, _, err := database_wiiu.GetPIDsByFriendRequestID(id)
	if err != nil {
		if err == database.ErrFriendRequestNotFound {
			return nil, nex.Errors.FPD.InvalidMessageID
		} else {
			globals.Logger.Critical(err.Error())
			return nil, nex.Errors.FPD.Unknown
		}
	}

	err = database_wiiu.SetUserBlocked(client.PID().LegacyValue(), senderPID, 0, 0)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return nil, nex.Errors.FPD.Unknown
	}

	info, err := utility.GetUserInfoByPID(senderPID)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return nil, nex.Errors.FPD.Unknown
	}

	date := nex.NewDateTime(0)
	date.FromTimestamp(time.Now())

	// Create a new blacklist principal for the client, as unlike AddBlacklist they don't send one to us.
	blacklistPrincipal := friends_wiiu_types.NewBlacklistedPrincipal()

	blacklistPrincipal.PrincipalBasicInfo = info
	blacklistPrincipal.GameKey = friends_wiiu_types.NewGameKey()
	blacklistPrincipal.BlackListedSince = date

	rmcResponseStream := nex.NewStreamOut(globals.SecureServer)

	rmcResponseStream.WriteStructure(blacklistPrincipal)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(globals.SecureServer, rmcResponseBody)
	rmcResponse.ProtocolID = friends_wiiu.ProtocolID
	rmcResponse.MethodID = friends_wiiu.MethodDenyFriendRequest
	rmcResponse.CallID = callID

	return rmcResponse, 0
}
