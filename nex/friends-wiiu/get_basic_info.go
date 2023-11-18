package nex_friends_wiiu

import (
	"github.com/PretendoNetwork/friends/globals"
	"github.com/PretendoNetwork/friends/utility"
	nex "github.com/PretendoNetwork/nex-go"
	friends_wiiu "github.com/PretendoNetwork/nex-protocols-go/friends-wiiu"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/friends-wiiu/types"
)

func GetBasicInfo(err error, packet nex.PacketInterface, callID uint32, pids []*nex.PID) (*nex.RMCMessage, uint32) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.Errors.FPD.InvalidArgument
	}

	infos := make([]*friends_wiiu_types.PrincipalBasicInfo, 0)

	for i := 0; i < len(pids); i++ {
		pid := pids[i]

		info, err := utility.GetUserInfoByPID(pid.LegacyValue())
		if err != nil {
			globals.Logger.Critical(err.Error())
			return nil, nex.Errors.FPD.Unknown
		}

		if info.PID.LegacyValue() != 0 {
			infos = append(infos, info)
		}
	}

	rmcResponseStream := nex.NewStreamOut(globals.SecureServer)

	nex.StreamWriteListStructure(rmcResponseStream, infos)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(rmcResponseBody)
	rmcResponse.ProtocolID = friends_wiiu.ProtocolID
	rmcResponse.MethodID = friends_wiiu.MethodGetBasicInfo
	rmcResponse.CallID = callID

	return rmcResponse, 0
}
