package nex_secure_connection

import (
	"net"
	"strconv"
	"time"

	database_3ds "github.com/PretendoNetwork/friends/database/3ds"
	database_wiiu "github.com/PretendoNetwork/friends/database/wiiu"
	"github.com/PretendoNetwork/friends/globals"
	"github.com/PretendoNetwork/friends/types"
	nex "github.com/PretendoNetwork/nex-go"
	secure_connection "github.com/PretendoNetwork/nex-protocols-go/secure-connection"
)

func RegisterEx(err error, packet nex.PacketInterface, callID uint32, stationUrls []*nex.StationURL, loginData *nex.DataHolder) (*nex.RMCMessage, uint32) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.Errors.Core.InvalidArgument
	}

	client := packet.Sender().(*nex.PRUDPClient)

	retval := nex.NewResultSuccess(nex.Errors.Core.Unknown)
	rmcResponseStream := nex.NewStreamOut(globals.SecureServer)

	// TODO - Validate loginData
	pid := client.PID().LegacyValue()

	user := types.NewConnectedUser()
	user.PID = pid
	user.Client = client

	lastOnline := nex.NewDateTime(0)
	lastOnline.FromTimestamp(time.Now())

	loginDataType := loginData.TypeName()
	switch loginDataType {
	case "NintendoLoginData":
		user.Platform = types.WUP // * Platform is Wii U

		err = database_wiiu.UpdateUserLastOnlineTime(pid, lastOnline)
		if err != nil {
			globals.Logger.Critical(err.Error())
			retval = nex.NewResultError(nex.Errors.Authentication.Unknown)
		}
	case "AccountExtraInfo":
		user.Platform = types.CTR // * Platform is 3DS

		err = database_3ds.UpdateUserLastOnlineTime(pid, lastOnline)
		if err != nil {
			globals.Logger.Critical(err.Error())
			retval = nex.NewResultError(nex.Errors.Authentication.Unknown)
		}
	default:
		globals.Logger.Errorf("Unknown loginData data type %s!", loginDataType)
		retval = nex.NewResultError(nex.Errors.Authentication.ValidationFailed)
	}

	if retval.IsSuccess() {
		globals.ConnectedUsers[pid] = user

		localStation := stationUrls[0]

		address := client.Address().(*net.UDPAddr)

		localStation.Fields.Set("address", address.IP.String())
		localStation.Fields.Set("port", strconv.Itoa(address.Port))

		localStationURL := localStation.EncodeToString()

		rmcResponseStream.WriteResult(retval)
		rmcResponseStream.WriteUInt32LE(globals.SecureServer.ConnectionIDCounter().Next())
		rmcResponseStream.WriteString(localStationURL)
	} else {
		rmcResponseStream.WriteResult(retval)
		rmcResponseStream.WriteUInt32LE(0)
		rmcResponseStream.WriteString("prudp:/")
	}

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(rmcResponseBody)
	rmcResponse.ProtocolID = secure_connection.ProtocolID
	rmcResponse.MethodID = secure_connection.MethodRegisterEx
	rmcResponse.CallID = callID

	return rmcResponse, 0
}
