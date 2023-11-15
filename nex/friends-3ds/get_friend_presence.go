package nex_friends_3ds

import (
	"github.com/PretendoNetwork/friends/globals"
	nex "github.com/PretendoNetwork/nex-go"
	friends_3ds "github.com/PretendoNetwork/nex-protocols-go/friends-3ds"
	friends_3ds_types "github.com/PretendoNetwork/nex-protocols-go/friends-3ds/types"
)

func GetFriendPresence(err error, packet nex.PacketInterface, callID uint32, pids []*nex.PID) (*nex.RMCMessage, uint32) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.Errors.FPD.Unknown
	}

	presenceList := make([]*friends_3ds_types.FriendPresence, 0)

	for i := 0; i < len(pids); i++ {
		connectedUser := globals.ConnectedUsers[pids[i].LegacyValue()]

		if connectedUser != nil && connectedUser.Presence != nil {
			friendPresence := friends_3ds_types.NewFriendPresence()
			friendPresence.PID = pids[i]
			friendPresence.Presence = globals.ConnectedUsers[pids[i].LegacyValue()].Presence

			presenceList = append(presenceList, friendPresence)
		}
	}

	rmcResponseStream := nex.NewStreamOut(globals.SecureServer)

	rmcResponseStream.WriteListStructure(presenceList)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(rmcResponseBody)
	rmcResponse.ProtocolID = friends_3ds.ProtocolID
	rmcResponse.MethodID = friends_3ds.MethodGetFriendPresence
	rmcResponse.CallID = callID

	return rmcResponse, 0
}
