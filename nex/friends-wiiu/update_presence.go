package nex_friends_wiiu

import (
	"github.com/PretendoNetwork/friends-secure/globals"
	notifications_wiiu "github.com/PretendoNetwork/friends-secure/notifications/wiiu"
	"github.com/PretendoNetwork/friends-secure/types"
	nex "github.com/PretendoNetwork/nex-go"
	friends_wiiu "github.com/PretendoNetwork/nex-protocols-go/friends/wiiu"
)

func UpdatePresence(err error, client *nex.Client, callID uint32, presence *friends_wiiu.NintendoPresenceV2) {
	pid := client.PID()

	presence.Online = true // Force online status. I have no idea why this is always false
	presence.PID = pid     // WHY IS THIS SET TO 0 BY DEFAULT??

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

	rmcResponse := nex.NewRMCResponse(friends_wiiu.ProtocolID, callID)
	rmcResponse.SetSuccess(friends_wiiu.MethodUpdatePresence, nil)

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewPacketV0(client, nil)

	responsePacket.SetVersion(0)
	responsePacket.SetSource(0xA1)
	responsePacket.SetDestination(0xAF)
	responsePacket.SetType(nex.DataPacket)
	responsePacket.SetPayload(rmcResponseBytes)

	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)

	globals.NEXServer.Send(responsePacket)
}
