package nex_friends_wiiu

import (
	"github.com/PretendoNetwork/friends/globals"
	notifications_wiiu "github.com/PretendoNetwork/friends/notifications/wiiu"
	friends_types "github.com/PretendoNetwork/friends/types"
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	friends_wiiu "github.com/PretendoNetwork/nex-protocols-go/v2/friends-wiiu"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/v2/friends-wiiu/types"
)

func UpdatePresence(err error, packet nex.PacketInterface, callID uint32, presence friends_wiiu_types.NintendoPresenceV2) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.FPD.InvalidArgument, "") // TODO - Add error message
	}

	connection := packet.Sender().(*nex.PRUDPConnection)

	pid := uint32(connection.PID())

	presence.Online = types.NewBool(true) // * Force online status. I have no idea why this is always false
	presence.PID = connection.PID()       // * WHY IS THIS SET TO 0 BY DEFAULT??

	connectedUser, ok := globals.ConnectedUsers.Get(pid)

	if !ok || connectedUser == nil {
		// TODO - Figure out why this is getting removed
		connectedUser = friends_types.NewConnectedUser()
		connectedUser.PID = pid
		connectedUser.Platform = friends_types.WUP
		connectedUser.Connection = connection
		// TODO - Find a clean way to create a NNAInfo?

		globals.ConnectedUsers.Set(pid, connectedUser)
	}

	connectedUser.PresenceV2 = presence.Copy().(friends_wiiu_types.NintendoPresenceV2)

	notifications_wiiu.SendPresenceUpdate(presence)

	rmcResponse := nex.NewRMCSuccess(globals.SecureEndpoint, nil)
	rmcResponse.ProtocolID = friends_wiiu.ProtocolID
	rmcResponse.MethodID = friends_wiiu.MethodUpdatePresence
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
