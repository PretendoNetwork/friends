package nex_friends_wiiu

import (
	database_wiiu "github.com/PretendoNetwork/friends/database/wiiu"
	"github.com/PretendoNetwork/friends/globals"
	nex "github.com/PretendoNetwork/nex-go"
	friends_wiiu "github.com/PretendoNetwork/nex-protocols-go/friends-wiiu"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/friends-wiiu/types"
)

func UpdateComment(err error, packet nex.PacketInterface, callID uint32, comment *friends_wiiu_types.Comment) (*nex.RMCMessage, uint32) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.Errors.FPD.InvalidArgument
	}

	client := packet.Sender().(*nex.PRUDPClient)

	changed, err := database_wiiu.UpdateUserComment(client.PID().LegacyValue(), comment.Contents)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return nil, nex.Errors.FPD.Unknown
	}

	rmcResponseStream := nex.NewStreamOut(globals.SecureServer)

	rmcResponseStream.WriteUInt64LE(changed)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(globals.SecureServer, rmcResponseBody)
	rmcResponse.ProtocolID = friends_wiiu.ProtocolID
	rmcResponse.MethodID = friends_wiiu.MethodUpdateComment
	rmcResponse.CallID = callID

	return rmcResponse, 0
}
