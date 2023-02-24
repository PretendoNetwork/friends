package friends_wiiu

import (
	"github.com/PretendoNetwork/friends-secure/globals"
	notifications_wiiu "github.com/PretendoNetwork/friends-secure/notifications/wiiu"
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func UpdatePresence(err error, client *nex.Client, callID uint32, presence *nexproto.NintendoPresenceV2) {
	pid := client.PID()

	presence.Online = true      // Force online status. I have no idea why this is always false
	presence.PID = client.PID() // WHY IS THIS SET TO 0 BY DEFAULT??

	globals.ConnectedUsers[pid].PresenceV2 = presence
	notifications_wiiu.SendPresenceUpdate(presence)

	rmcResponse := nex.NewRMCResponse(nexproto.FriendsWiiUProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.FriendsWiiUMethodUpdatePresence, nil)

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
