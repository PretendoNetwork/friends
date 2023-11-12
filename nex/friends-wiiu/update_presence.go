package nex_friends_wiiu

import (
	"github.com/PretendoNetwork/friends/globals"
	notifications_wiiu "github.com/PretendoNetwork/friends/notifications/wiiu"
	"github.com/PretendoNetwork/friends/types"
	nex "github.com/PretendoNetwork/nex-go"
	friends_wiiu "github.com/PretendoNetwork/nex-protocols-go/friends-wiiu"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/friends-wiiu/types"
)

func UpdatePresence(err error, packet nex.PacketInterface, callID uint32, presence *friends_wiiu_types.NintendoPresenceV2) uint32 {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nex.Errors.FPD.InvalidArgument
	}

	client := packet.Sender().(*nex.PRUDPClient)

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

	rmcResponse := nex.NewRMCSuccess(nil)
	rmcResponse.ProtocolID = friends_wiiu.ProtocolID
	rmcResponse.MethodID = friends_wiiu.MethodUpdatePresence
	rmcResponse.CallID = callID

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewPRUDPPacketV0(client, nil)

	responsePacket.SetType(nex.DataPacket)
	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)
	responsePacket.SetSourceStreamType(packet.(nex.PRUDPPacketInterface).DestinationStreamType())
	responsePacket.SetSourcePort(packet.(nex.PRUDPPacketInterface).DestinationPort())
	responsePacket.SetDestinationStreamType(packet.(nex.PRUDPPacketInterface).SourceStreamType())
	responsePacket.SetDestinationPort(packet.(nex.PRUDPPacketInterface).SourcePort())
	responsePacket.SetPayload(rmcResponseBytes)

	globals.SecureServer.Send(responsePacket)

	return 0
}
