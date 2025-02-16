package nex_friends_3ds

import (
	"database/sql"

	"github.com/PretendoNetwork/friends/database"
	database_3ds "github.com/PretendoNetwork/friends/database/3ds"
	"github.com/PretendoNetwork/friends/globals"
	notifications_3ds "github.com/PretendoNetwork/friends/notifications/3ds"
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	friends_3ds "github.com/PretendoNetwork/nex-protocols-go/v2/friends-3ds"
	friends_3ds_types "github.com/PretendoNetwork/nex-protocols-go/v2/friends-3ds/types"
)

func SyncFriend(err error, packet nex.PacketInterface, callID uint32, lfc types.UInt64, pids types.List[types.PID], lfcList types.List[types.UInt64]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.FPD.InvalidArgument, "") // TODO - Add error message
	}

	connection := packet.Sender().(*nex.PRUDPConnection)

	friendRelationships, err := database_3ds.GetUserFriends(uint32(connection.PID()))
	if err != nil && err != sql.ErrNoRows {
		globals.Logger.Critical(err.Error())
		return nil, nex.NewError(nex.ResultCodes.FPD.Unknown, "") // TODO - Add error message
	}

	for _, relationship := range friendRelationships {
		var hasPID bool
		for _, pid := range pids {
			if pid == relationship.PID {
				hasPID = true
				break
			}
		}

		if !hasPID {
			err := database_3ds.RemoveFriendship(uint32(connection.PID()), uint32(relationship.PID))
			if err != nil && err != database.ErrFriendshipNotFound {
				globals.Logger.Critical(err.Error())
				return nil, nex.NewError(nex.ResultCodes.FPD.Unknown, "") // TODO - Add error message
			}
		}
	}

	// TODO - Not needed?
	relationships := friendRelationships.Copy().(types.List[friends_3ds_types.FriendRelationship])

	for _, pid := range pids {
		if !isPIDInRelationships(relationships, uint32(pid)) {
			relationship, err := database_3ds.SaveFriendship(uint32(connection.PID()), uint32(pid))
			if err != nil {
				globals.Logger.Critical(err.Error())
				return nil, nex.NewError(nex.ResultCodes.FPD.Unknown, "") // TODO - Add error message
			}

			relationships = append(relationships, relationship)

			// * Alert the other side, in case they weren't able to get our presence data
			connectedUser, ok := globals.ConnectedUsers.Get(uint32(pid))
			if ok && connectedUser != nil {
				go notifications_3ds.SendFriendshipCompleted(connectedUser.Connection, connection.PID())
			}
		}
	}

	// TODO - Not needed?
	syncedRelationships := relationships.Copy().(types.List[friends_3ds_types.FriendRelationship])

	rmcResponseStream := nex.NewByteStreamOut(globals.SecureEndpoint.LibraryVersions(), globals.SecureEndpoint.ByteStreamSettings())

	syncedRelationships.WriteTo(rmcResponseStream)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(globals.SecureEndpoint, rmcResponseBody)
	rmcResponse.ProtocolID = friends_3ds.ProtocolID
	rmcResponse.MethodID = friends_3ds.MethodSyncFriend
	rmcResponse.CallID = callID

	return rmcResponse, nil
}

func isPIDInRelationships(relationships []friends_3ds_types.FriendRelationship, pid uint32) bool {
	for i := range relationships {
		if uint32(relationships[i].PID) == pid {
			return true
		}
	}

	return false
}
