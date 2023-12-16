package nex_friends_wiiu

import (
	database_wiiu "github.com/PretendoNetwork/friends/database/wiiu"
	"github.com/PretendoNetwork/friends/globals"
	nex "github.com/PretendoNetwork/nex-go"
	friends_wiiu "github.com/PretendoNetwork/nex-protocols-go/friends-wiiu"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/friends-wiiu/types"
)

func UpdatePreference(err error, packet nex.PacketInterface, callID uint32, principalPreference *friends_wiiu_types.PrincipalPreference) (*nex.RMCMessage, uint32) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.Errors.FPD.InvalidArgument
	}

	client := packet.Sender().(*nex.PRUDPClient)

	err = database_wiiu.UpdateUserPrincipalPreference(client.PID().LegacyValue(), principalPreference)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return nil, nex.Errors.FPD.Unknown
	}

	rmcResponse := nex.NewRMCSuccess(globals.SecureServer, nil)
	rmcResponse.ProtocolID = friends_wiiu.ProtocolID
	rmcResponse.MethodID = friends_wiiu.MethodUpdatePreference
	rmcResponse.CallID = callID

	return rmcResponse, 0
}
