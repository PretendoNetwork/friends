package friends_3ds

import (
	database_3ds "github.com/PretendoNetwork/friends-secure/database/3ds"
	"github.com/PretendoNetwork/friends-secure/globals"
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func GetAllFriends(err error, client *nex.Client, callID uint32) {
	friendRelationships := database_3ds.GetUserFriends(client.PID())

	rmcResponseStream := nex.NewStreamOut(globals.NEXServer)

	rmcResponseStream.WriteListStructure(friendRelationships)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(nexproto.Friends3DSProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.Friends3DSMethodGetAllFriends, rmcResponseBody)

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewPacketV0(client, nil)

	responsePacket.SetVersion(0)
	responsePacket.SetSource(0xA1)
	responsePacket.SetDestination(0xAF)
	responsePacket.SetType(nex.DataPacket)
	responsePacket.SetPayload(rmcResponseBytes)

	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)

	globals.NEXServer.Send(responsePacket)
}