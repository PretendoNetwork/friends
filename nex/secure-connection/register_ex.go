package nex_secure_connection

import (
	"net"

	database_3ds "github.com/PretendoNetwork/friends/database/3ds"
	database_wiiu "github.com/PretendoNetwork/friends/database/wiiu"
	"github.com/PretendoNetwork/friends/globals"
	friends_types "github.com/PretendoNetwork/friends/types"
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	secure_connection "github.com/PretendoNetwork/nex-protocols-go/v2/secure-connection"
)

func RegisterEx(err error, packet nex.PacketInterface, callID uint32, vecMyURLs *types.List[*types.StationURL], hCustomData *types.AnyDataHolder) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "")
	}

	connection := packet.Sender().(*nex.PRUDPConnection)

	retval := types.NewQResultSuccess(nex.ResultCodes.Core.Unknown)

	// TODO - Validate loginData
	pid := connection.PID().LegacyValue()

	user := friends_types.NewConnectedUser()
	user.PID = pid
	user.Connection = connection

	lastOnline := types.NewDateTime(0).Now()
	loginDataType := hCustomData.TypeName.Value

	switch loginDataType {
	case "NintendoLoginData":
		user.Platform = friends_types.WUP // * Platform is Wii U

		err = database_wiiu.UpdateUserLastOnlineTime(pid, lastOnline)
		if err != nil {
			globals.Logger.Critical(err.Error())
			retval = types.NewQResultError(nex.ResultCodes.Authentication.Unknown)
		}
	case "AccountExtraInfo":
		user.Platform = friends_types.CTR // * Platform is 3DS

		err = database_3ds.UpdateUserLastOnlineTime(pid, lastOnline)
		if err != nil {
			globals.Logger.Critical(err.Error())
			retval = types.NewQResultError(nex.ResultCodes.Authentication.Unknown)
		}
	default:
		globals.Logger.Errorf("Unknown loginData data type %s!", loginDataType)
		retval = types.NewQResultError(nex.ResultCodes.Authentication.ValidationFailed)
	}

	pidConnectionID := types.NewPrimitiveU32(0)
	urlPublic := types.NewString("")

	if retval.IsSuccess() {
		globals.ConnectedUsers[pid] = user

		localStation, _ := vecMyURLs.Get(0)

		address := connection.Address().(*net.UDPAddr)

		localStation.SetAddress(address.IP.String())
		localStation.SetPortNumber(uint16(address.Port))

		localStationURL := localStation.EncodeToString()

		pidConnectionID = types.NewPrimitiveU32(connection.ID)
		urlPublic = types.NewString(localStationURL)
	}

	rmcResponseStream := nex.NewByteStreamOut(globals.SecureEndpoint.LibraryVersions(), globals.SecureEndpoint.ByteStreamSettings())

	retval.WriteTo(rmcResponseStream)
	pidConnectionID.WriteTo(rmcResponseStream)
	urlPublic.WriteTo(rmcResponseStream)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(globals.SecureEndpoint, rmcResponseBody)
	rmcResponse.ProtocolID = secure_connection.ProtocolID
	rmcResponse.MethodID = secure_connection.MethodRegisterEx
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
