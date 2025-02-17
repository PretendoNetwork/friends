package nex_friends_3ds

import (
	"github.com/PretendoNetwork/friends/database"
	database_3ds "github.com/PretendoNetwork/friends/database/3ds"
	"github.com/PretendoNetwork/friends/globals"
	notifications_3ds "github.com/PretendoNetwork/friends/notifications/3ds"
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	friends_3ds "github.com/PretendoNetwork/nex-protocols-go/v2/friends-3ds"
)

func RemoveFriendByPrincipalID(err error, packet nex.PacketInterface, callID uint32, pid types.PID) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.FPD.InvalidArgument, "") // TODO - Add error message
	}

	connection := packet.Sender().(*nex.PRUDPConnection)

	err = database_3ds.RemoveFriendship(uint32(connection.PID()), uint32(pid))
	if err != nil {
		if err == database.ErrFriendshipNotFound {
			// * Official servers don't actually check this, but
			// * we'll do it ourselves
			return nil, nex.NewError(nex.ResultCodes.FPD.NotFriend, "") // TODO - Add error message
		} else {
			globals.Logger.Critical(err.Error())
			return nil, nex.NewError(nex.ResultCodes.FPD.Unknown, "") // TODO - Add error message
		}
	}

	go notifications_3ds.SendUserWentOffline(connection, pid)

	rmcResponse := nex.NewRMCSuccess(globals.SecureEndpoint, nil)
	rmcResponse.ProtocolID = friends_3ds.ProtocolID
	rmcResponse.MethodID = friends_3ds.MethodRemoveFriendByPrincipalID
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
