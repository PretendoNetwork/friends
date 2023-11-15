package nex_friends_3ds

import (
	database_3ds "github.com/PretendoNetwork/friends/database/3ds"
	"github.com/PretendoNetwork/friends/globals"
	notifications_3ds "github.com/PretendoNetwork/friends/notifications/3ds"
	nex "github.com/PretendoNetwork/nex-go"
	friends_3ds "github.com/PretendoNetwork/nex-protocols-go/friends-3ds"
)

func AddFriendshipByPrincipalID(err error, packet nex.PacketInterface, callID uint32, lfc uint64, pid *nex.PID) (*nex.RMCMessage, uint32) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.Errors.FPD.InvalidArgument
	}

	client := packet.Sender().(*nex.PRUDPClient)

	friendRelationship, err := database_3ds.SaveFriendship(client.PID().LegacyValue(), pid.LegacyValue())
	if err != nil {
		globals.Logger.Critical(err.Error())
		return nil, nex.Errors.FPD.Unknown
	}

	connectedUser := globals.ConnectedUsers[pid.LegacyValue()]
	if connectedUser != nil {
		go notifications_3ds.SendFriendshipCompleted(connectedUser.Client, pid.LegacyValue(), client.PID())
	}

	rmcResponseStream := nex.NewStreamOut(globals.SecureServer)

	rmcResponseStream.WriteStructure(friendRelationship)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(rmcResponseBody)
	rmcResponse.ProtocolID = friends_3ds.ProtocolID
	rmcResponse.MethodID = friends_3ds.MethodAddFriendByPrincipalID
	rmcResponse.CallID = callID

	return rmcResponse, 0
}
