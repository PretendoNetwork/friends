package nex_friends_wiiu

import (
	database_wiiu "github.com/PretendoNetwork/friends/database/wiiu"
	"github.com/PretendoNetwork/friends/globals"
	nex "github.com/PretendoNetwork/nex-go"
	friends_wiiu "github.com/PretendoNetwork/nex-protocols-go/friends-wiiu"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/friends-wiiu/types"
)

func GetRequestBlockSettings(err error, packet nex.PacketInterface, callID uint32, pids []uint32) (*nex.RMCMessage, uint32) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.Errors.FPD.InvalidArgument
	}

	client := packet.Sender().(*nex.PRUDPClient)

	settings := make([]*friends_wiiu_types.PrincipalRequestBlockSetting, 0)

	// TODO:
	// Improve this. Use less database_wiiu.reads
	for i := 0; i < len(pids); i++ {
		requestedPID := pids[i]

		setting := friends_wiiu_types.NewPrincipalRequestBlockSetting()
		setting.PID = requestedPID
		isBlocked, err := database_wiiu.IsFriendRequestBlocked(client.PID().LegacyValue(), requestedPID)
		if err != nil {
			globals.Logger.Critical(err.Error())
			return nil, nex.Errors.Core.Unknown
		}

		setting.IsBlocked = isBlocked

		settings = append(settings, setting)
	}

	rmcResponseStream := nex.NewStreamOut(globals.SecureServer)

	nex.StreamWriteListStructure(rmcResponseStream, settings)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(globals.SecureServer, rmcResponseBody)
	rmcResponse.ProtocolID = friends_wiiu.ProtocolID
	rmcResponse.MethodID = friends_wiiu.MethodGetRequestBlockSettings
	rmcResponse.CallID = callID

	return rmcResponse, 0
}
