package nex_friends_wiiu

import (
	"github.com/PretendoNetwork/friends/globals"
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	friends_wiiu "github.com/PretendoNetwork/nex-protocols-go/v2/friends-wiiu"
)

func CheckSettingStatus(err error, packet nex.PacketInterface, callID uint32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.FPD.Unknown, "") // TODO - Add error message
	}

	status := types.NewUInt8(0xFF) // TODO - What is this??

	rmcResponseStream := nex.NewByteStreamOut(globals.SecureEndpoint.LibraryVersions(), globals.SecureEndpoint.ByteStreamSettings())

	status.WriteTo(rmcResponseStream)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(globals.SecureEndpoint, rmcResponseBody)
	rmcResponse.ProtocolID = friends_wiiu.ProtocolID
	rmcResponse.MethodID = friends_wiiu.MethodCheckSettingStatus
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
