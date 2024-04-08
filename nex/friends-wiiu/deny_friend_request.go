package nex_friends_wiiu

import (
	"time"

	"github.com/PretendoNetwork/friends/database"
	database_wiiu "github.com/PretendoNetwork/friends/database/wiiu"
	"github.com/PretendoNetwork/friends/globals"
	"github.com/PretendoNetwork/friends/utility"
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	friends_wiiu "github.com/PretendoNetwork/nex-protocols-go/v2/friends-wiiu"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/v2/friends-wiiu/types"
)

func DenyFriendRequest(err error, packet nex.PacketInterface, callID uint32, id *types.PrimitiveU64) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.FPD.InvalidArgument, "") // TODO - Add error message
	}

	connection := packet.Sender().(*nex.PRUDPConnection)

	err = database_wiiu.SetFriendRequestDenied(id.Value)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return nil, nex.NewError(nex.ResultCodes.FPD.Unknown, "") // TODO - Add error message
	}

	senderPID, _, err := database_wiiu.GetPIDsByFriendRequestID(id.Value)
	if err != nil {
		if err == database.ErrFriendRequestNotFound {
			return nil, nex.NewError(nex.ResultCodes.FPD.InvalidMessageID, "") // TODO - Add error message
		} else {
			globals.Logger.Critical(err.Error())
			return nil, nex.NewError(nex.ResultCodes.FPD.Unknown, "") // TODO - Add error message
		}
	}

	err = database_wiiu.SetUserBlocked(connection.PID().LegacyValue(), senderPID, 0, 0)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return nil, nex.NewError(nex.ResultCodes.FPD.Unknown, "") // TODO - Add error message
	}

	info, err := utility.GetUserInfoByPID(senderPID)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return nil, nex.NewError(nex.ResultCodes.FPD.Unknown, "") // TODO - Add error message
	}

	date := types.NewDateTime(0)
	date.FromTimestamp(time.Now())

	// Create a new blacklist principal for the connection, as unlike AddBlacklist they don't send one to us.
	blacklistPrincipal := friends_wiiu_types.NewBlacklistedPrincipal()

	blacklistPrincipal.PrincipalBasicInfo = info
	blacklistPrincipal.GameKey = friends_wiiu_types.NewGameKey()
	blacklistPrincipal.BlackListedSince = date

	rmcResponseStream := nex.NewByteStreamOut(globals.SecureEndpoint.LibraryVersions(), globals.SecureEndpoint.ByteStreamSettings())

	blacklistPrincipal.WriteTo(rmcResponseStream)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(globals.SecureEndpoint, rmcResponseBody)
	rmcResponse.ProtocolID = friends_wiiu.ProtocolID
	rmcResponse.MethodID = friends_wiiu.MethodDenyFriendRequest
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
