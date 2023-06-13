package nex_friends_3ds

import (
	"github.com/PretendoNetwork/friends-secure/globals"
	notifications_3ds "github.com/PretendoNetwork/friends-secure/notifications/3ds"
	"github.com/PretendoNetwork/friends-secure/types"
	nex "github.com/PretendoNetwork/nex-go"
	friends_3ds "github.com/PretendoNetwork/nex-protocols-go/friends/3ds"
)

func UpdatePresence(err error, client *nex.Client, callID uint32, presence *friends_3ds.NintendoPresence, showGame bool) {
	currentPresence := presence

	// Send an entirely empty status, with every flag set to update
	if !showGame {
		currentPresence = friends_3ds.NewNintendoPresence()
		currentPresence.GameKey = friends_3ds.NewGameKey()
		currentPresence.ChangedFlags = 0xFFFFFFFF // All flags
	}

	go notifications_3ds.SendPresenceUpdate(client, currentPresence)

	pid := client.PID()

	if globals.ConnectedUsers[pid] == nil {
		// TODO - Figure out why this is getting removed
		connectedUser := types.NewConnectedUser()
		connectedUser.PID = pid
		connectedUser.Platform = types.CTR
		connectedUser.Client = client

		globals.ConnectedUsers[pid] = connectedUser
	}

	globals.ConnectedUsers[pid].Presence = currentPresence

	rmcResponse := nex.NewRMCResponse(friends_3ds.ProtocolID, callID)
	rmcResponse.SetSuccess(friends_3ds.MethodUpdatePresence, nil)

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
