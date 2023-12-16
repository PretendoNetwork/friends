package nex_friends_3ds

import (
	"github.com/PretendoNetwork/friends/globals"
	notifications_3ds "github.com/PretendoNetwork/friends/notifications/3ds"
	"github.com/PretendoNetwork/friends/types"
	nex "github.com/PretendoNetwork/nex-go"
	friends_3ds "github.com/PretendoNetwork/nex-protocols-go/friends-3ds"
	friends_3ds_types "github.com/PretendoNetwork/nex-protocols-go/friends-3ds/types"
)

func UpdatePresence(err error, packet nex.PacketInterface, callID uint32, presence *friends_3ds_types.NintendoPresence, showGame bool) (*nex.RMCMessage, uint32) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.Errors.FPD.InvalidArgument
	}

	client := packet.Sender().(*nex.PRUDPClient)

	currentPresence := presence

	// Send an entirely empty status, with every flag set to update
	if !showGame {
		currentPresence = friends_3ds_types.NewNintendoPresence()
		currentPresence.GameKey = friends_3ds_types.NewGameKey()
		currentPresence.ChangedFlags = 0xFFFFFFFF // All flags
	}

	go notifications_3ds.SendPresenceUpdate(client, currentPresence)

	pid := client.PID().LegacyValue()

	if globals.ConnectedUsers[pid] == nil {
		// TODO - Figure out why this is getting removed
		connectedUser := types.NewConnectedUser()
		connectedUser.PID = pid
		connectedUser.Platform = types.CTR
		connectedUser.Client = client

		globals.ConnectedUsers[pid] = connectedUser
	}

	globals.ConnectedUsers[pid].Presence = currentPresence

	rmcResponse := nex.NewRMCSuccess(globals.SecureServer, nil)
	rmcResponse.ProtocolID = friends_3ds.ProtocolID
	rmcResponse.MethodID = friends_3ds.MethodUpdatePresence
	rmcResponse.CallID = callID

	return rmcResponse, 0
}
