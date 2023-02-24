package friends_3ds

import (
	"github.com/PretendoNetwork/friends-secure/globals"
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func GetFriendPresence(err error, client *nex.Client, callID uint32, pids []uint32) {
	presenceList := make([]*nexproto.FriendPresence, 0)

	for i := 0; i < len(pids); i++ {
		connectedUser := globals.ConnectedUsers[pids[i]]

		if connectedUser != nil && connectedUser.Presence != nil {
			friendPresence := nexproto.NewFriendPresence()
			friendPresence.PID = pids[i]
			friendPresence.Presence = globals.ConnectedUsers[pids[i]].Presence

			presenceList = append(presenceList, friendPresence)
		}
	}

	rmcResponseStream := nex.NewStreamOut(globals.NEXServer)

	rmcResponseStream.WriteListStructure(presenceList)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(nexproto.Friends3DSProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.Friends3DSMethodGetFriendPresence, rmcResponseBody)

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
