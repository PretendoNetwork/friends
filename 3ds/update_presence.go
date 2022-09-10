package friends_3ds

import (
	"github.com/PretendoNetwork/friends-secure/globals"
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func UpdatePresence(err error, client *nex.Client, callID uint32, presence *nexproto.NintendoPresence, unknown bool) {
	// TODO: Do something with this

	rmcResponse := nex.NewRMCResponse(nexproto.Friends3DSProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.Friends3DSMethodUpdatePresence, nil)

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
