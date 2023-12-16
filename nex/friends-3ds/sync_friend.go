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
)

func SyncFriend(err error, packet nex.PacketInterface, callID uint32, lfc uint64, pids []*nex.PID, lfcList []uint64) (*nex.RMCMessage, uint32) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.Errors.FPD.InvalidArgument
	}

	client := packet.Sender().(*nex.PRUDPClient)

	friendRelationships, err := database_3ds.GetUserFriends(client.PID().LegacyValue())
	if err != nil && err != sql.ErrNoRows {
		globals.Logger.Critical(err.Error())
		return nil, nex.Errors.FPD.Unknown
	}

	for i := 0; i < len(friendRelationships); i++ {
		var hasPID bool
		for _, pidInput := range pids {
			if pidInput.Equals(friendRelationships[i].PID) {
				hasPID = true
				break
			}
		}

		if !hasPID {
			err := database_3ds.RemoveFriendship(client.PID().LegacyValue(), friendRelationships[i].PID.LegacyValue())
			if err != nil && err != database.ErrFriendshipNotFound {
				globals.Logger.Critical(err.Error())
				return nil, nex.Errors.FPD.Unknown
			}
		}
	}

	for i := 0; i < len(pids); i++ {
		if !isPIDInRelationships(friendRelationships, pids[i].LegacyValue()) {
			friendRelationship, err := database_3ds.SaveFriendship(client.PID().LegacyValue(), pids[i].LegacyValue())
			if err != nil {
				globals.Logger.Critical(err.Error())
				return nil, nex.Errors.FPD.Unknown
			}

			friendRelationships = append(friendRelationships, friendRelationship)

			// Alert the other side, in case they weren't able to get our presence data
			connectedUser := globals.ConnectedUsers[pids[i].LegacyValue()]
			if connectedUser != nil {
				go notifications_3ds.SendFriendshipCompleted(connectedUser.Client, pids[i].LegacyValue(), client.PID())
			}
		}
	}

	rmcResponseStream := nex.NewStreamOut(globals.SecureServer)

	nex.StreamWriteListStructure(rmcResponseStream, friendRelationships)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(globals.SecureServer, rmcResponseBody)
	rmcResponse.ProtocolID = friends_3ds.ProtocolID
	rmcResponse.MethodID = friends_3ds.MethodSyncFriend
	rmcResponse.CallID = callID

	return rmcResponse, 0
}

func isPIDInRelationships(relationships []*friends_3ds_types.FriendRelationship, pid uint32) bool {
	for i := range relationships {
		if relationships[i].PID.LegacyValue() == pid {
			return true
		}
	}

	return false
}
