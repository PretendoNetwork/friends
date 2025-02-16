package nex_friends_wiiu

import (
	database_wiiu "github.com/PretendoNetwork/friends/database/wiiu"
	"github.com/PretendoNetwork/friends/globals"
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	friends_wiiu "github.com/PretendoNetwork/nex-protocols-go/v2/friends-wiiu"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/v2/friends-wiiu/types"
)

func GetBasicInfo(err error, packet nex.PacketInterface, callID uint32, pids types.List[types.PID]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.FPD.InvalidArgument, "") // TODO - Add error message
	}

	infos := types.NewList[friends_wiiu_types.PrincipalBasicInfo]()

	for _, pid := range pids {
		info, err := database_wiiu.GetUserPrincipalBasicInfo(uint32(pid))
		if err != nil {
			globals.Logger.Critical(err.Error())
			return nil, nex.NewError(nex.ResultCodes.FPD.Unknown, "") // TODO - Add error message
		}

		if info.PID != 0 {
			infos = append(infos, info)
		}
	}

	rmcResponseStream := nex.NewByteStreamOut(globals.SecureEndpoint.LibraryVersions(), globals.SecureEndpoint.ByteStreamSettings())

	infos.WriteTo(rmcResponseStream)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(globals.SecureEndpoint, rmcResponseBody)
	rmcResponse.ProtocolID = friends_wiiu.ProtocolID
	rmcResponse.MethodID = friends_wiiu.MethodGetBasicInfo
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
