package main

import (
	"strconv"

	"github.com/PretendoNetwork/friends-secure/globals"
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func registerEx(err error, client *nex.Client, callID uint32, stationUrls []*nex.StationURL, loginData nexproto.NintendoLoginData) {
	// TODO: Validate loginData

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
