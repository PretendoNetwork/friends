package nex_friends_3ds

import (
	"database/sql"

	database_3ds "github.com/PretendoNetwork/friends/database/3ds"
	"github.com/PretendoNetwork/friends/globals"
	nex "github.com/PretendoNetwork/nex-go"
	friends_3ds "github.com/PretendoNetwork/nex-protocols-go/friends-3ds"
)

func GetFriendPersistentInfo(err error, client *nex.Client, callID uint32, pids []uint32) uint32 {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nex.Errors.FPD.Unknown
	}

	infoList, err := database_3ds.GetFriendPersistentInfos(client.PID(), pids)
	if err != nil && err != sql.ErrNoRows {
		globals.Logger.Critical(err.Error())
		return nex.Errors.FPD.Unknown
	}

	rmcResponseStream := nex.NewStreamOut(globals.SecureServer)

	rmcResponseStream.WriteListStructure(infoList)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(friends_3ds.ProtocolID, callID)
	rmcResponse.SetSuccess(friends_3ds.MethodGetFriendPersistentInfo, rmcResponseBody)

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewPacketV0(client, nil)

	responsePacket.SetVersion(0)
	responsePacket.SetSource(0xA1)
	responsePacket.SetDestination(0xAF)
	responsePacket.SetType(nex.DataPacket)
	responsePacket.SetPayload(rmcResponseBytes)

	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)

	globals.SecureServer.Send(responsePacket)

	return 0
}
