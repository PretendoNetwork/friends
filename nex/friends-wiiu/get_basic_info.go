package nex_friends_wiiu

import (
	"github.com/PretendoNetwork/friends/globals"
	"github.com/PretendoNetwork/friends/utility"
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	friends_wiiu "github.com/PretendoNetwork/nex-protocols-go/v2/friends-wiiu"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/v2/friends-wiiu/types"
)

func GetBasicInfo(err error, packet nex.PacketInterface, callID uint32, pids *types.List[*types.PID]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.FPD.InvalidArgument, "") // TODO - Add error message
	}

	infos := types.NewList[*friends_wiiu_types.PrincipalBasicInfo]()
	infos.Type = friends_wiiu_types.NewPrincipalBasicInfo()

	if pids.Each(func(i int, pid *types.PID) bool {
		info, err := utility.GetUserInfoByPID(pid.LegacyValue())
		if err != nil {
			globals.Logger.Critical(err.Error())
			return true
		}

		if info.PID.LegacyValue() != 0 {
			infos.Append(info)
		}

		return false
	}) {
		return nil, nex.NewError(nex.ResultCodes.FPD.Unknown, "") // TODO - Add error message
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
