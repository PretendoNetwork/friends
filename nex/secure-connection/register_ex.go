package nex_secure_connection

import (
	"strconv"
	"time"

	database_3ds "github.com/PretendoNetwork/friends/database/3ds"
	database_wiiu "github.com/PretendoNetwork/friends/database/wiiu"
	"github.com/PretendoNetwork/friends/globals"
	"github.com/PretendoNetwork/friends/types"
	nex "github.com/PretendoNetwork/nex-go"
	secure_connection "github.com/PretendoNetwork/nex-protocols-go/secure-connection"
)

func RegisterEx(err error, client *nex.Client, callID uint32, stationUrls []*nex.StationURL, loginData *nex.DataHolder) uint32 {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nex.Errors.Core.InvalidArgument
	}

	retval := nex.NewResultSuccess(nex.Errors.Core.Unknown)
	rmcResponseStream := nex.NewStreamOut(globals.SecureServer)

	// TODO: Validate loginData
	pid := client.PID()
	user := globals.ConnectedUsers[pid]
	lastOnline := nex.NewDateTime(0)
	lastOnline.FromTimestamp(time.Now())

	loginDataType := loginData.TypeName()
	switch loginDataType {
	case "NintendoLoginData":
		user.Platform = types.WUP // Platform is Wii U

		database_wiiu.UpdateUserLastOnlineTime(pid, lastOnline)
	case "AccountExtraInfo":
		user.Platform = types.CTR // Platform is 3DS

		database_3ds.UpdateUserLastOnlineTime(pid, lastOnline)
	default:
		globals.Logger.Errorf("Unknown loginData data type %s!", loginDataType)
		retval = nex.NewResultError(nex.Errors.Authentication.ValidationFailed)
	}

	if retval.IsSuccess() {
		localStation := stationUrls[0]

		address := client.Address().IP.String()
		port := strconv.Itoa(client.Address().Port)

		localStation.SetAddress(address)
		localStation.SetPort(port)

		localStationURL := localStation.EncodeToString()

		rmcResponseStream.WriteResult(retval)
		rmcResponseStream.WriteUInt32LE(globals.SecureServer.ConnectionIDCounter().Increment())
		rmcResponseStream.WriteString(localStationURL)
	} else {
		rmcResponseStream.WriteResult(retval)
		rmcResponseStream.WriteUInt32LE(0)
		rmcResponseStream.WriteString("prudp:/")
	}

	rmcResponseBody := rmcResponseStream.Bytes()

	// Build response packet
	rmcResponse := nex.NewRMCResponse(secure_connection.ProtocolID, callID)
	rmcResponse.SetSuccess(secure_connection.MethodRegisterEx, rmcResponseBody)

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
