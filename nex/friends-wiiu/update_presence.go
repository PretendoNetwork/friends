package nex_friends_wiiu

import (
	"github.com/PretendoNetwork/friends/globals"
	notifications_wiiu "github.com/PretendoNetwork/friends/notifications/wiiu"
	"github.com/PretendoNetwork/friends/types"
	nex "github.com/PretendoNetwork/nex-go"
	friends_wiiu "github.com/PretendoNetwork/nex-protocols-go/friends-wiiu"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/friends-wiiu/types"
)

func UpdatePresence(err error, packet nex.PacketInterface, callID uint32, presence *friends_wiiu_types.NintendoPresenceV2) (*nex.RMCMessage, uint32) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.Errors.FPD.InvalidArgument
	}

	client := packet.Sender().(*nex.PRUDPClient)

	pid := client.PID().LegacyValue()

	presence.Online = true      // * Force online status. I have no idea why this is always false
	presence.PID = client.PID() // * WHY IS THIS SET TO 0 BY DEFAULT??

	if globals.ConnectedUsers[pid] == nil {
		// TODO - Figure out why this is getting removed
		connectedUser := types.NewConnectedUser()
		connectedUser.PID = pid
		connectedUser.Platform = types.WUP
		connectedUser.Client = client
		// TODO - Find a clean way to create a NNAInfo?

		globals.ConnectedUsers[pid] = connectedUser
	}

	globals.ConnectedUsers[pid].PresenceV2 = presence

	notifications_wiiu.SendPresenceUpdate(presence)

	rmcResponse := nex.NewRMCSuccess(globals.SecureServer, nil)
	rmcResponse.ProtocolID = friends_wiiu.ProtocolID
	rmcResponse.MethodID = friends_wiiu.MethodUpdatePresence
	rmcResponse.CallID = callID

	return rmcResponse, 0
}
