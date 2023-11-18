package nex_friends_3ds

import (
	"database/sql"

	database_3ds "github.com/PretendoNetwork/friends/database/3ds"
	"github.com/PretendoNetwork/friends/globals"
	nex "github.com/PretendoNetwork/nex-go"
	friends_3ds "github.com/PretendoNetwork/nex-protocols-go/friends-3ds"
)

func GetFriendPersistentInfo(err error, packet nex.PacketInterface, callID uint32, pids []*nex.PID) (*nex.RMCMessage, uint32) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.Errors.FPD.Unknown
	}

	client := packet.Sender().(*nex.PRUDPClient)

	friendPIDs := make([]uint32, len(pids))

	for _, pid := range pids {
		friendPIDs = append(friendPIDs, pid.LegacyValue())
	}

	infoList, err := database_3ds.GetFriendPersistentInfos(client.PID().LegacyValue(), friendPIDs)
	if err != nil && err != sql.ErrNoRows {
		globals.Logger.Critical(err.Error())
		return nil, nex.Errors.FPD.Unknown
	}

	rmcResponseStream := nex.NewStreamOut(globals.SecureServer)

	nex.StreamWriteListStructure(rmcResponseStream, infoList)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(rmcResponseBody)
	rmcResponse.ProtocolID = friends_3ds.ProtocolID
	rmcResponse.MethodID = friends_3ds.MethodGetFriendPersistentInfo
	rmcResponse.CallID = callID

	return rmcResponse, 0
}
