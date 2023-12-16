package nex_friends_wiiu

import (
	"github.com/PretendoNetwork/friends/globals"
	nex "github.com/PretendoNetwork/nex-go"
	friends_wiiu "github.com/PretendoNetwork/nex-protocols-go/friends-wiiu"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/friends-wiiu/types"
)

func DeletePersistentNotification(err error, packet nex.PacketInterface, callID uint32, notifications []*friends_wiiu_types.PersistentNotification) (*nex.RMCMessage, uint32) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.Errors.FPD.InvalidArgument
	}

	// TODO - Do something here

	rmcResponse := nex.NewRMCSuccess(globals.SecureServer, nil)
	rmcResponse.ProtocolID = friends_wiiu.ProtocolID
	rmcResponse.MethodID = friends_wiiu.MethodDeletePersistentNotification
	rmcResponse.CallID = callID

	return rmcResponse, 0
}
