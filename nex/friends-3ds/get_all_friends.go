package nex_friends_3ds

import (
	"database/sql"

	database_3ds "github.com/PretendoNetwork/friends/database/3ds"
	"github.com/PretendoNetwork/friends/globals"
	nex "github.com/PretendoNetwork/nex-go/v2"
	friends_3ds "github.com/PretendoNetwork/nex-protocols-go/v2/friends-3ds"
)

func GetAllFriends(err error, packet nex.PacketInterface, callID uint32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.FPD.Unknown, "") // TODO - Add error message
	}

	connection := packet.Sender().(*nex.PRUDPConnection)

	friendRelationships, err := database_3ds.GetUserFriends(uint32(connection.PID()))
	if err != nil && err != sql.ErrNoRows {
		globals.Logger.Critical(err.Error())
		return nil, nex.NewError(nex.ResultCodes.FPD.Unknown, "") // TODO - Add error message
	}

	rmcResponseStream := nex.NewByteStreamOut(globals.SecureEndpoint.LibraryVersions(), globals.SecureEndpoint.ByteStreamSettings())

	friendRelationships.WriteTo(rmcResponseStream)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(globals.SecureEndpoint, rmcResponseBody)
	rmcResponse.ProtocolID = friends_3ds.ProtocolID
	rmcResponse.MethodID = friends_3ds.MethodGetAllFriends
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
