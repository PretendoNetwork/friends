package nex_friends_wiiu

import (
	"github.com/PretendoNetwork/friends/database"
	database_wiiu "github.com/PretendoNetwork/friends/database/wiiu"
	"github.com/PretendoNetwork/friends/globals"
	notifications_wiiu "github.com/PretendoNetwork/friends/notifications/wiiu"
	nex "github.com/PretendoNetwork/nex-go"
	friends_wiiu "github.com/PretendoNetwork/nex-protocols-go/friends-wiiu"
)

func CancelFriendRequest(err error, packet nex.PacketInterface, callID uint32, id uint64) (*nex.RMCMessage, uint32) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.Errors.FPD.InvalidArgument
	}

	client := packet.Sender().(*nex.PRUDPClient)

	pid, err := database_wiiu.DeleteFriendRequestAndReturnFriendPID(id)
	if err != nil {
		if err == database.ErrFriendRequestNotFound {
			return nil, nex.Errors.FPD.InvalidMessageID
		} else {
			globals.Logger.Critical(err.Error())
			return nil, nex.Errors.FPD.Unknown
		}
	}

	connectedUser := globals.ConnectedUsers[pid]
	if connectedUser != nil {
		// * This may send the friend removed notification, but they are the same.
		go notifications_wiiu.SendFriendshipRemoved(connectedUser.Client, client.PID())
	}

	rmcResponse := nex.NewRMCSuccess(nil)
	rmcResponse.ProtocolID = friends_wiiu.ProtocolID
	rmcResponse.MethodID = friends_wiiu.MethodCancelFriendRequest
	rmcResponse.CallID = callID

	return rmcResponse, 0
}
