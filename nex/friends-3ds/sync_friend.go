package nex_friends_3ds

import (
	"database/sql"

	"github.com/PretendoNetwork/friends/database"
	database_3ds "github.com/PretendoNetwork/friends/database/3ds"
	"github.com/PretendoNetwork/friends/globals"
	notifications_3ds "github.com/PretendoNetwork/friends/notifications/3ds"
	nex "github.com/PretendoNetwork/nex-go"
	friends_3ds "github.com/PretendoNetwork/nex-protocols-go/friends-3ds"
	friends_3ds_types "github.com/PretendoNetwork/nex-protocols-go/friends-3ds/types"
	"golang.org/x/exp/slices"
)

func SyncFriend(err error, packet nex.PacketInterface, callID uint32, lfc uint64, pids []uint32, lfcList []uint64) uint32 {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nex.Errors.FPD.InvalidArgument
	}

	client := packet.Sender().(*nex.PRUDPClient)

	friendRelationships, err := database_3ds.GetUserFriends(client.PID().LegacyValue())
	if err != nil && err != sql.ErrNoRows {
		globals.Logger.Critical(err.Error())
		return nex.Errors.FPD.Unknown
	}

	for i := 0; i < len(friendRelationships); i++ {
		if !slices.Contains(pids, friendRelationships[i].PID) {
			err := database_3ds.RemoveFriendship(client.PID().LegacyValue(), friendRelationships[i].PID)
			if err != nil && err != database.ErrFriendshipNotFound {
				globals.Logger.Critical(err.Error())
				return nex.Errors.FPD.Unknown
			}
		}
	}

	for i := 0; i < len(pids); i++ {
		if !isPIDInRelationships(friendRelationships, pids[i]) {
			friendRelationship, err := database_3ds.SaveFriendship(client.PID().LegacyValue(), pids[i])
			if err != nil {
				globals.Logger.Critical(err.Error())
				return nex.Errors.FPD.Unknown
			}

			friendRelationships = append(friendRelationships, friendRelationship)

			// Alert the other side, in case they weren't able to get our presence data
			connectedUser := globals.ConnectedUsers[pids[i]]
			if connectedUser != nil {
				go notifications_3ds.SendFriendshipCompleted(connectedUser.Client, pids[i], client.PID().LegacyValue())
			}
		}
	}

	rmcResponseStream := nex.NewStreamOut(globals.SecureServer)

	rmcResponseStream.WriteListStructure(friendRelationships)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(rmcResponseBody)
	rmcResponse.ProtocolID = friends_3ds.ProtocolID
	rmcResponse.MethodID = friends_3ds.MethodSyncFriend
	rmcResponse.CallID = callID

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewPRUDPPacketV0(client, nil)

	responsePacket.SetType(nex.DataPacket)
	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)
	responsePacket.SetSourceStreamType(packet.(nex.PRUDPPacketInterface).DestinationStreamType())
	responsePacket.SetSourcePort(packet.(nex.PRUDPPacketInterface).DestinationPort())
	responsePacket.SetDestinationStreamType(packet.(nex.PRUDPPacketInterface).SourceStreamType())
	responsePacket.SetDestinationPort(packet.(nex.PRUDPPacketInterface).SourcePort())
	responsePacket.SetPayload(rmcResponseBytes)

	globals.SecureServer.Send(responsePacket)

	return 0
}

func isPIDInRelationships(relationships []*friends_3ds_types.FriendRelationship, pid uint32) bool {
	for i := range relationships {
		if relationships[i].PID == pid {
			return true
		}
	}
	return false
}
