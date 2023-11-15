package nex_friends_wiiu

import (
	"github.com/PretendoNetwork/friends/database"
	database_wiiu "github.com/PretendoNetwork/friends/database/wiiu"
	"github.com/PretendoNetwork/friends/globals"
	"github.com/PretendoNetwork/friends/utility"
	nex "github.com/PretendoNetwork/nex-go"
	friends_wiiu "github.com/PretendoNetwork/nex-protocols-go/friends-wiiu"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/friends-wiiu/types"
)

func AddBlacklist(err error, packet nex.PacketInterface, callID uint32, blacklistPrincipal *friends_wiiu_types.BlacklistedPrincipal) (*nex.RMCMessage, uint32) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.Errors.FPD.InvalidArgument
	}

	client := packet.Sender().(*nex.PRUDPClient)

	currentBlacklistPrincipal := blacklistPrincipal

	senderPID := currentBlacklistPrincipal.PrincipalBasicInfo.PID
	titleID := currentBlacklistPrincipal.GameKey.TitleID
	titleVersion := currentBlacklistPrincipal.GameKey.TitleVersion

	userInfo, err := utility.GetUserInfoByPID(currentBlacklistPrincipal.PrincipalBasicInfo.PID.LegacyValue())
	if err != nil {
		if err == database.ErrPIDNotFound {
			return nil, nex.Errors.FPD.InvalidPrincipalID // TODO - Not sure if this is the correct error.
		} else {
			globals.Logger.Critical(err.Error())
			return nil, nex.Errors.FPD.Unknown
		}
	}

	currentBlacklistPrincipal.PrincipalBasicInfo = userInfo
	currentBlacklistPrincipal.BlackListedSince = nex.NewDateTime(0).Now()

	err = database_wiiu.SetUserBlocked(client.PID().LegacyValue(), senderPID.LegacyValue(), titleID, titleVersion)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return nil, nex.Errors.FPD.Unknown
	}

	rmcResponseStream := nex.NewStreamOut(globals.SecureServer)

	rmcResponseStream.WriteStructure(blacklistPrincipal)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(rmcResponseBody)
	rmcResponse.ProtocolID = friends_wiiu.ProtocolID
	rmcResponse.MethodID = friends_wiiu.MethodAddBlackList
	rmcResponse.CallID = callID

	return rmcResponse, 0
}
