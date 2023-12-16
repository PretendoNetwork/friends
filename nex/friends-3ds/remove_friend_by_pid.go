package nex_friends_3ds

import (
	"github.com/PretendoNetwork/friends/database"
	database_3ds "github.com/PretendoNetwork/friends/database/3ds"
	"github.com/PretendoNetwork/friends/globals"
	notifications_3ds "github.com/PretendoNetwork/friends/notifications/3ds"
	nex "github.com/PretendoNetwork/nex-go"
	friends_3ds "github.com/PretendoNetwork/nex-protocols-go/friends-3ds"
)

func RemoveFriendByPrincipalID(err error, packet nex.PacketInterface, callID uint32, pid *nex.PID) (*nex.RMCMessage, uint32) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.Errors.FPD.InvalidArgument
	}

	client := packet.Sender().(*nex.PRUDPClient)

	err = database_3ds.RemoveFriendship(client.PID().LegacyValue(), pid.LegacyValue())
	if err != nil {
		if err == database.ErrFriendshipNotFound {
			// * Official servers don't actually check this, but
			// * we'll do it ourselves
			return nil, nex.Errors.FPD.NotFriend
		} else {
			globals.Logger.Critical(err.Error())
			return nil, nex.Errors.FPD.Unknown
		}
	}

	go notifications_3ds.SendUserWentOffline(client, pid)

	rmcResponse := nex.NewRMCSuccess(globals.SecureServer, nil)
	rmcResponse.ProtocolID = friends_3ds.ProtocolID
	rmcResponse.MethodID = friends_3ds.MethodRemoveFriendByPrincipalID
	rmcResponse.CallID = callID

	return rmcResponse, 0
}
