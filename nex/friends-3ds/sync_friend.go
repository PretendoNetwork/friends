package nex_friends_3ds

import (
	database_3ds "github.com/PretendoNetwork/friends-secure/database/3ds"
	"github.com/PretendoNetwork/friends-secure/globals"
	notifications_3ds "github.com/PretendoNetwork/friends-secure/notifications/3ds"
	nex "github.com/PretendoNetwork/nex-go"
	friends_3ds "github.com/PretendoNetwork/nex-protocols-go/friends/3ds"
	"golang.org/x/exp/slices"
)

func SyncFriend(err error, client *nex.Client, callID uint32, lfc uint64, pids []uint32, lfcList []uint64) {
	friendRelationships := database_3ds.GetUserFriends(client.PID())

	for i := 0; i < len(friendRelationships); i++ {
		if !slices.Contains(pids, friendRelationships[i].PID) {
			database_3ds.RemoveFriendship(client.PID(), friendRelationships[i].PID)
		}
	}

	for i := 0; i < len(pids); i++ {
		if !isPIDInRelationships(friendRelationships, pids[i]) {
			friendRelationship := database_3ds.SaveFriendship(client.PID(), pids[i])

			friendRelationships = append(friendRelationships, friendRelationship)

			// Alert the other side, in case they weren't able to get our presence data
			connectedUser := globals.ConnectedUsers[pids[i]]
			if connectedUser != nil {
				go notifications_3ds.SendFriendshipCompleted(connectedUser.Client, pids[i], client.PID())
			}
		}
	}

	rmcResponseStream := nex.NewStreamOut(globals.NEXServer)

	rmcResponseStream.WriteListStructure(friendRelationships)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(friends_3ds.ProtocolID, callID)
	rmcResponse.SetSuccess(friends_3ds.MethodSyncFriend, rmcResponseBody)

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

func isPIDInRelationships(relationships []*friends_3ds.FriendRelationship, pid uint32) bool {
	for i := range relationships {
		if relationships[i].PID == pid {
			return true
		}
	}
	return false
}
