package nex_friends_3ds

import (
	database_3ds "github.com/PretendoNetwork/friends/database/3ds"
	"github.com/PretendoNetwork/friends/globals"
	nex "github.com/PretendoNetwork/nex-go/v2"
	friends_3ds "github.com/PretendoNetwork/nex-protocols-go/v2/friends-3ds"
	friends_3ds_types "github.com/PretendoNetwork/nex-protocols-go/v2/friends-3ds/types"
)

func UpdateProfile(err error, packet nex.PacketInterface, callID uint32, profileData friends_3ds_types.MyProfile) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.FPD.InvalidArgument, "") // TODO - Add error message
	}

	connection := packet.Sender().(*nex.PRUDPConnection)

	err = database_3ds.UpdateUserProfile(uint32(connection.PID()), profileData)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return nil, nex.NewError(nex.ResultCodes.FPD.Unknown, "") // TODO - Add error message
	}

	rmcResponse := nex.NewRMCSuccess(globals.SecureEndpoint, nil)
	rmcResponse.ProtocolID = friends_3ds.ProtocolID
	rmcResponse.MethodID = friends_3ds.MethodUpdateProfile
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
