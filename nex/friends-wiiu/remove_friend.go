package nex_friends_wiiu

import (
	"github.com/PretendoNetwork/friends/database"
	database_wiiu "github.com/PretendoNetwork/friends/database/wiiu"
	"github.com/PretendoNetwork/friends/globals"
	notifications_wiiu "github.com/PretendoNetwork/friends/notifications/wiiu"
	nex "github.com/PretendoNetwork/nex-go"
	friends_wiiu "github.com/PretendoNetwork/nex-protocols-go/friends-wiiu"
)

func RemoveFriend(err error, packet nex.PacketInterface, callID uint32, pid *nex.PID) (*nex.RMCMessage, uint32) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.Errors.FPD.InvalidArgument
	}

	client := packet.Sender().(*nex.PRUDPClient)

	err = database_wiiu.RemoveFriendship(client.PID().LegacyValue(), pid.LegacyValue())
	if err != nil {
		if err == database.ErrFriendshipNotFound {
			return nil, nex.Errors.FPD.NotInMyFriendList
		} else {
			globals.Logger.Critical(err.Error())
			return nil, nex.Errors.FPD.Unknown
		}
	}

	connectedUser := globals.ConnectedUsers[pid.LegacyValue()]
	if connectedUser != nil {
		go notifications_wiiu.SendFriendshipRemoved(connectedUser.Client, pid)
	}

	rmcResponse := nex.NewRMCSuccess(nil)
	rmcResponse.ProtocolID = friends_wiiu.ProtocolID
	rmcResponse.MethodID = friends_wiiu.MethodRemoveFriend
	rmcResponse.CallID = callID

	return rmcResponse, 0
}
