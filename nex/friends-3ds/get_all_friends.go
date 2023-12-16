package nex_friends_3ds

import (
	"database/sql"

	database_3ds "github.com/PretendoNetwork/friends/database/3ds"
	"github.com/PretendoNetwork/friends/globals"
	nex "github.com/PretendoNetwork/nex-go"
	friends_3ds "github.com/PretendoNetwork/nex-protocols-go/friends-3ds"
)

func GetAllFriends(err error, packet nex.PacketInterface, callID uint32) (*nex.RMCMessage, uint32) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.Errors.FPD.Unknown
	}

	client := packet.Sender().(*nex.PRUDPClient)

	friendRelationships, err := database_3ds.GetUserFriends(client.PID().LegacyValue())
	if err != nil && err != sql.ErrNoRows {
		globals.Logger.Critical(err.Error())
		return nil, nex.Errors.FPD.Unknown
	}

	rmcResponseStream := nex.NewStreamOut(globals.SecureServer)

	nex.StreamWriteListStructure(rmcResponseStream, friendRelationships)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(globals.SecureServer, rmcResponseBody)
	rmcResponse.ProtocolID = friends_3ds.ProtocolID
	rmcResponse.MethodID = friends_3ds.MethodGetAllFriends
	rmcResponse.CallID = callID

	return rmcResponse, 0
}
