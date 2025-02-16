package nex_friends_3ds

import (
	"database/sql"

	database_3ds "github.com/PretendoNetwork/friends/database/3ds"
	"github.com/PretendoNetwork/friends/globals"
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	friends_3ds "github.com/PretendoNetwork/nex-protocols-go/v2/friends-3ds"
)

func GetFriendPersistentInfo(err error, packet nex.PacketInterface, callID uint32, pidList types.List[types.PID]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.FPD.Unknown, "") // TODO - Add error message
	}

	connection := packet.Sender().(*nex.PRUDPConnection)

	friendPIDs := make([]uint32, len(pidList))

	for _, pid := range pidList {
		friendPIDs = append(friendPIDs, uint32(pid))
	}

	infoList, err := database_3ds.GetFriendPersistentInfos(uint32(connection.PID()), friendPIDs)
	if err != nil && err != sql.ErrNoRows {
		globals.Logger.Critical(err.Error())
		return nil, nex.NewError(nex.ResultCodes.FPD.Unknown, "") // TODO - Add error message
	}

	rmcResponseStream := nex.NewByteStreamOut(globals.SecureEndpoint.LibraryVersions(), globals.SecureEndpoint.ByteStreamSettings())

	infoList.WriteTo(rmcResponseStream)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(globals.SecureEndpoint, rmcResponseBody)
	rmcResponse.ProtocolID = friends_3ds.ProtocolID
	rmcResponse.MethodID = friends_3ds.MethodGetFriendPersistentInfo
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
