package nex_friends_wiiu

import (
	database_wiiu "github.com/PretendoNetwork/friends/database/wiiu"
	"github.com/PretendoNetwork/friends/globals"
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	friends_wiiu "github.com/PretendoNetwork/nex-protocols-go/v2/friends-wiiu"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/v2/friends-wiiu/types"
)

func UpdateComment(err error, packet nex.PacketInterface, callID uint32, comment friends_wiiu_types.Comment) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.FPD.InvalidArgument, "") // TODO - Add error message
	}

	connection := packet.Sender().(*nex.PRUDPConnection)

	changed, err := database_wiiu.UpdateUserComment(uint32(connection.PID()), string(comment.Contents))
	if err != nil {
		globals.Logger.Critical(err.Error())
		return nil, nex.NewError(nex.ResultCodes.FPD.Unknown, "") // TODO - Add error message
	}

	rmcResponseStream := nex.NewByteStreamOut(globals.SecureEndpoint.LibraryVersions(), globals.SecureEndpoint.ByteStreamSettings())

	types.NewUInt64(changed).WriteTo(rmcResponseStream) // TODO - This is ugly

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(globals.SecureEndpoint, rmcResponseBody)
	rmcResponse.ProtocolID = friends_wiiu.ProtocolID
	rmcResponse.MethodID = friends_wiiu.MethodUpdateComment
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
