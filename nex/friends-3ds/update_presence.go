package nex_friends_3ds

import (
	"github.com/PretendoNetwork/friends/globals"
	notifications_3ds "github.com/PretendoNetwork/friends/notifications/3ds"
	friends_types "github.com/PretendoNetwork/friends/types"
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	friends_3ds "github.com/PretendoNetwork/nex-protocols-go/v2/friends-3ds"
	friends_3ds_types "github.com/PretendoNetwork/nex-protocols-go/v2/friends-3ds/types"
)

func UpdatePresence(err error, packet nex.PacketInterface, callID uint32, presence *friends_3ds_types.NintendoPresence, showGame *types.PrimitiveBool) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.FPD.InvalidArgument, "") // TODO - Add error message
	}

	connection := packet.Sender().(*nex.PRUDPConnection)

	currentPresence := presence.Copy().(*friends_3ds_types.NintendoPresence)

	// Send an entirely empty status, with every flag set to update
	if !showGame.Value {
		currentPresence = friends_3ds_types.NewNintendoPresence()
		currentPresence.GameKey = friends_3ds_types.NewGameKey()
		currentPresence.ChangedFlags = types.NewPrimitiveU32(0xFFFFFFFF) // * All flags
	}

	go notifications_3ds.SendPresenceUpdate(connection, currentPresence)

	pid := connection.PID().LegacyValue()
	connectedUser, ok := globals.ConnectedUsers.Get(pid)

	if !ok || connectedUser == nil {
		// TODO - Figure out why this is getting removed
		connectedUser = friends_types.NewConnectedUser()
		connectedUser.PID = pid
		connectedUser.Platform = friends_types.CTR
		connectedUser.Connection = connection

		globals.ConnectedUsers.Set(pid, connectedUser)
	}

	connectedUser.Presence = currentPresence

	rmcResponse := nex.NewRMCSuccess(globals.SecureEndpoint, nil)
	rmcResponse.ProtocolID = friends_3ds.ProtocolID
	rmcResponse.MethodID = friends_3ds.MethodUpdatePresence
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
