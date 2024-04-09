package nex_friends_3ds

import (
	"github.com/PretendoNetwork/friends/globals"
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	friends_3ds "github.com/PretendoNetwork/nex-protocols-go/v2/friends-3ds"
	friends_3ds_types "github.com/PretendoNetwork/nex-protocols-go/v2/friends-3ds/types"
)

func GetFriendPresence(err error, packet nex.PacketInterface, callID uint32, pidList *types.List[*types.PID]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.FPD.Unknown, "") // TODO - Add error message
	}

	presenceList := types.NewList[*friends_3ds_types.FriendPresence]()
	presenceList.Type = friends_3ds_types.NewFriendPresence()

	pidList.Each(func(i int, pid *types.PID) bool {
		connectedUser, ok := globals.ConnectedUsers.Get(pid.LegacyValue())

		if ok && connectedUser != nil && connectedUser.Presence != nil {
			friendPresence := friends_3ds_types.NewFriendPresence()
			friendPresence.PID = pid.Copy().(*types.PID)
			friendPresence.Presence = connectedUser.Presence

			presenceList.Append(friendPresence)
		}

		return false
	})

	rmcResponseStream := nex.NewByteStreamOut(globals.SecureEndpoint.LibraryVersions(), globals.SecureEndpoint.ByteStreamSettings())

	presenceList.WriteTo(rmcResponseStream)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(globals.SecureEndpoint, rmcResponseBody)
	rmcResponse.ProtocolID = friends_3ds.ProtocolID
	rmcResponse.MethodID = friends_3ds.MethodGetFriendPresence
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
