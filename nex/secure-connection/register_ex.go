package nex_secure_connection

import (
	"strconv"
	"time"

	database_3ds "github.com/PretendoNetwork/friends-secure/database/3ds"
	database_wiiu "github.com/PretendoNetwork/friends-secure/database/wiiu"
	"github.com/PretendoNetwork/friends-secure/globals"
	"github.com/PretendoNetwork/friends-secure/types"
	nex "github.com/PretendoNetwork/nex-go"
	secure_connection "github.com/PretendoNetwork/nex-protocols-go/secure-connection"
)

func RegisterEx(err error, client *nex.Client, callID uint32, stationUrls []*nex.StationURL, loginData *nex.DataHolder) uint32 {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nex.Errors.Core.InvalidArgument
	}

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
		return nex.Errors.Authentication.TokenParseError
	}

	localStation := stationUrls[0]

	address := client.Address().IP.String()
	port := strconv.Itoa(client.Address().Port)

	localStation.SetAddress(address)
	localStation.SetPort(port)

	localStationURL := localStation.EncodeToString()

	rmcResponseStream := nex.NewStreamOut(globals.SecureServer)

	retval := nex.NewResultSuccess(nex.Errors.Core.Unknown)
	rmcResponseStream.WriteResult(retval)
	rmcResponseStream.WriteUInt32LE(globals.SecureServer.ConnectionIDCounter().Increment())
	rmcResponseStream.WriteString(localStationURL)

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
