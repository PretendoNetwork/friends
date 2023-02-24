package main

import (
	"strconv"
	"time"

	database_3ds "github.com/PretendoNetwork/friends-secure/database/3ds"
	database_wiiu "github.com/PretendoNetwork/friends-secure/database/wiiu"
	"github.com/PretendoNetwork/friends-secure/globals"
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func registerEx(err error, client *nex.Client, callID uint32, stationUrls []*nex.StationURL, loginData *nex.DataHolder) {
	// TODO: Validate loginData

	pid := client.PID()
	user := globals.ConnectedUsers[pid]
	lastOnline := nex.NewDateTime(0)
	lastOnline.FromTimestamp(time.Now())

	if loginData.TypeName() == "NintendoLoginData" {
		user.Platform = 1 // Platform is Wii U

		database_wiiu.UpdateUserLastOnlineTime(pid, lastOnline)
	} else if loginData.TypeName() == "AccountExtraInfo" {
		user.Platform = 2 // Platform is 3DS

		database_3ds.UpdateUserLastOnlineTime(pid, lastOnline)
	}

	localStation := stationUrls[0]

	address := client.Address().IP.String()
	port := strconv.Itoa(client.Address().Port)

	localStation.SetAddress(address)
	localStation.SetPort(port)

	localStationURL := localStation.EncodeToString()

	rmcResponseStream := nex.NewStreamOut(globals.NEXServer)

	rmcResponseStream.WriteUInt32LE(0x10001) // Success
	rmcResponseStream.WriteUInt32LE(globals.NEXServer.ConnectionIDCounter().Increment())
	rmcResponseStream.WriteString(localStationURL)

	rmcResponseBody := rmcResponseStream.Bytes()

	// Build response packet
	rmcResponse := nex.NewRMCResponse(nexproto.SecureProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.SecureMethodRegisterEx, rmcResponseBody)

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewPacketV0(client, nil)

	responsePacket.SetVersion(0)
	responsePacket.SetSource(0xA1)
	responsePacket.SetDestination(0xAF)
	responsePacket.SetType(nex.DataPacket)
	responsePacket.SetPayload(rmcResponseBytes)

	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)

	globals.NEXServer.Send(responsePacket)
}
