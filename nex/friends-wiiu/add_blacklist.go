package nex_friends_wiiu

import (
	"github.com/PretendoNetwork/friends/database"
	database_wiiu "github.com/PretendoNetwork/friends/database/wiiu"
	"github.com/PretendoNetwork/friends/globals"
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	friends_wiiu "github.com/PretendoNetwork/nex-protocols-go/v2/friends-wiiu"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/v2/friends-wiiu/types"
)

func AddBlackList(err error, packet nex.PacketInterface, callID uint32, blacklistPrincipal friends_wiiu_types.BlacklistedPrincipal) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.FPD.InvalidArgument, "") // TODO - Add error message
	}

	connection := packet.Sender().(*nex.PRUDPConnection)

	currentBlacklistPrincipal := blacklistPrincipal

	senderPID := currentBlacklistPrincipal.PrincipalBasicInfo.PID
	titleID := currentBlacklistPrincipal.GameKey.TitleID
	titleVersion := currentBlacklistPrincipal.GameKey.TitleVersion

	userInfo, err := database_wiiu.GetUserPrincipalBasicInfo(uint32(currentBlacklistPrincipal.PrincipalBasicInfo.PID))
	if err != nil {
		if err == database.ErrPIDNotFound {
			// TODO - Not sure if this is the correct error.
			return nil, nex.NewError(nex.ResultCodes.FPD.InvalidPrincipalID, "") // TODO - Add error message
		} else {
			globals.Logger.Critical(err.Error())
			return nil, nex.NewError(nex.ResultCodes.FPD.Unknown, "") // TODO - Add error message
		}
	}

	currentBlacklistPrincipal.PrincipalBasicInfo = userInfo
	currentBlacklistPrincipal.BlackListedSince = types.NewDateTime(0).Now()

	err = database_wiiu.SetUserBlocked(uint32(connection.PID()), uint32(senderPID), uint64(titleID), uint16(titleVersion))
	if err != nil {
		globals.Logger.Critical(err.Error())
		return nil, nex.NewError(nex.ResultCodes.FPD.Unknown, "") // TODO - Add error message
	}

	rmcResponseStream := nex.NewByteStreamOut(globals.SecureEndpoint.LibraryVersions(), globals.SecureEndpoint.ByteStreamSettings())

	blacklistPrincipal.WriteTo(rmcResponseStream)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(globals.SecureEndpoint, rmcResponseBody)
	rmcResponse.ProtocolID = friends_wiiu.ProtocolID
	rmcResponse.MethodID = friends_wiiu.MethodAddBlackList
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
