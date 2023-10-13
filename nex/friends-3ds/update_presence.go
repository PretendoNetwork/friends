package nex_friends_3ds

import (
	"github.com/PretendoNetwork/friends/globals"
	notifications_3ds "github.com/PretendoNetwork/friends/notifications/3ds"
	"github.com/PretendoNetwork/friends/types"
	nex "github.com/PretendoNetwork/nex-go"
	friends_3ds "github.com/PretendoNetwork/nex-protocols-go/friends-3ds"
	friends_3ds_types "github.com/PretendoNetwork/nex-protocols-go/friends-3ds/types"
)

func UpdatePresence(err error, client *nex.Client, callID uint32, presence *friends_3ds_types.NintendoPresence, showGame bool) uint32 {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nex.Errors.FPD.InvalidArgument
	}

	currentPresence := presence

	// Send an entirely empty status, with every flag set to update
	if !showGame {
		currentPresence = friends_3ds_types.NewNintendoPresence()
		currentPresence.GameKey = friends_3ds_types.NewGameKey()
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

	globals.SecureServer.Send(responsePacket)

	return 0
}
