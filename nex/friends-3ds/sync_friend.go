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

func SyncFriend(err error, packet nex.PacketInterface, callID uint32, lfc *types.PrimitiveU64, pids *types.List[*types.PID], lfcList *types.List[*types.PrimitiveU64]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.FPD.InvalidArgument, "") // TODO - Add error message
	}

	connection := packet.Sender().(*nex.PRUDPConnection)

	friendRelationships, err := database_3ds.GetUserFriends(connection.PID().LegacyValue())
	if err != nil && err != sql.ErrNoRows {
		globals.Logger.Critical(err.Error())
		return nil, nex.NewError(nex.ResultCodes.FPD.Unknown, "") // TODO - Add error message
	}

	if friendRelationships.Each(func(i int, relationship *friends_3ds_types.FriendRelationship) bool {
		var hasPID bool
		pids.Each(func(i int, pid *types.PID) bool {
			if pid.Equals(relationship.PID) {
				hasPID = true
				return true
			}

			return false
		})

		if !hasPID {
			err := database_3ds.RemoveFriendship(connection.PID().LegacyValue(), relationship.PID.LegacyValue())
			if err != nil && err != database.ErrFriendshipNotFound {
				globals.Logger.Critical(err.Error())
				return true
			}
		}

		return false
	}) {
		return nil, nex.NewError(nex.ResultCodes.FPD.Unknown, "") // TODO - Add error message
	}

	relationships := friendRelationships.Slice()

	if pids.Each(func(i int, pid *types.PID) bool {
		if !isPIDInRelationships(relationships, pid.LegacyValue()) {
			relationship, err := database_3ds.SaveFriendship(connection.PID().LegacyValue(), pid.LegacyValue())
			if err != nil {
				globals.Logger.Critical(err.Error())
				return true
			}

			relationships = append(relationships, relationship)

			// * Alert the other side, in case they weren't able to get our presence data
			connectedUser := globals.ConnectedUsers[pid.LegacyValue()]
			if connectedUser != nil {
				go notifications_3ds.SendFriendshipCompleted(connectedUser.Connection, pid.LegacyValue(), connection.PID())
			}
		}

		return false
	}) {
		return nil, nex.NewError(nex.ResultCodes.FPD.Unknown, "") // TODO - Add error message
	}

	syncedRelationships := types.NewList[*friends_3ds_types.FriendRelationship]()
	syncedRelationships.Type = friends_3ds_types.NewFriendRelationship()
	syncedRelationships.SetFromData(relationships)

	rmcResponseStream := nex.NewByteStreamOut(globals.SecureEndpoint.LibraryVersions(), globals.SecureEndpoint.ByteStreamSettings())

	syncedRelationships.WriteTo(rmcResponseStream)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(globals.SecureEndpoint, rmcResponseBody)
	rmcResponse.ProtocolID = friends_3ds.ProtocolID
	rmcResponse.MethodID = friends_3ds.MethodSyncFriend
	rmcResponse.CallID = callID

	return rmcResponse, nil
}

func isPIDInRelationships(relationships []*friends_3ds_types.FriendRelationship, pid uint32) bool {
	for i := range relationships {
		if relationships[i].PID.LegacyValue() == pid {
			return true
		}
	}

	return false
}
