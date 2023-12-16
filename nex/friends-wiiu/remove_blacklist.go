package nex_friends_wiiu

import (
	"github.com/PretendoNetwork/friends/database"
	database_wiiu "github.com/PretendoNetwork/friends/database/wiiu"
	"github.com/PretendoNetwork/friends/globals"
	nex "github.com/PretendoNetwork/nex-go"
	friends_wiiu "github.com/PretendoNetwork/nex-protocols-go/friends-wiiu"
)

func RemoveBlacklist(err error, packet nex.PacketInterface, callID uint32, blockedPID *nex.PID) (*nex.RMCMessage, uint32) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.Errors.FPD.InvalidArgument
	}

	client := packet.Sender().(*nex.PRUDPClient)

	err = database_wiiu.UnsetUserBlocked(client.PID().LegacyValue(), blockedPID.LegacyValue())
	if err != nil {
		if err == database.ErrPIDNotFound {
			return nil, nex.Errors.FPD.NotInMyBlacklist
		} else {
			globals.Logger.Critical(err.Error())
			return nil, nex.Errors.FPD.Unknown
		}
	}

	rmcResponse := nex.NewRMCSuccess(globals.SecureServer, nil)
	rmcResponse.ProtocolID = friends_wiiu.ProtocolID
	rmcResponse.MethodID = friends_wiiu.MethodRemoveBlackList
	rmcResponse.CallID = callID

	return rmcResponse, 0
}
