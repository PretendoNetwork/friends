package main

import (
	"time"

	"github.com/PretendoNetwork/friends-secure/database"
	"github.com/PretendoNetwork/friends-secure/globals"
	nex "github.com/PretendoNetwork/nex-go"
)

func connect(packet *nex.PacketV0) {
	packet.Sender().SetClientConnectionSignature(packet.ConnectionSignature())

	payload := packet.Payload()
	stream := nex.NewStreamIn(payload, globals.NEXServer)

	ticketData, _ := stream.ReadBuffer()
	requestData, _ := stream.ReadBuffer()

	serverKey := nex.DeriveKerberosKey(2, []byte(globals.NEXServer.KerberosPassword()))

	// TODO: use random key from auth server
	ticketDataEncryption := nex.NewKerberosEncryption(serverKey)
	decryptedTicketData := ticketDataEncryption.Decrypt(ticketData)
	ticketDataStream := nex.NewStreamIn(decryptedTicketData, globals.NEXServer)

	_ = ticketDataStream.ReadUInt64LE() // expiration time
	_ = ticketDataStream.ReadUInt32LE() // User PID
	sessionKey := ticketDataStream.ReadBytesNext(16)

	requestDataEncryption := nex.NewKerberosEncryption(sessionKey)
	decryptedRequestData := requestDataEncryption.Decrypt(requestData)
	requestDataStream := nex.NewStreamIn(decryptedRequestData, globals.NEXServer)

	userPID := requestDataStream.ReadUInt32LE() // User PID

	_ = requestDataStream.ReadUInt32LE() //CID of secure server station url
	responseCheck := requestDataStream.ReadUInt32LE()

	responseValueStream := nex.NewStreamOut(globals.NEXServer)
	responseValueStream.WriteUInt32LE(responseCheck + 1)

	responseValueBufferStream := nex.NewStreamOut(globals.NEXServer)
	responseValueBufferStream.WriteBuffer(responseValueStream.Bytes())

	packet.Sender().UpdateRC4Key(sessionKey)

	globals.NEXServer.AcknowledgePacket(packet, responseValueBufferStream.Bytes())

	packet.Sender().SetPID(userPID)

	connectedUser := globals.NewConnectedUser()
	connectedUser.PID = packet.Sender().PID()
	connectedUser.Client = packet.Sender()
	globals.ConnectedUsers[userPID] = connectedUser

	lastOnline := nex.NewDateTime(0)
	lastOnline.FromTimestamp(time.Now())

	database.UpdateUserLastOnlineTime(packet.Sender().PID(), lastOnline)
}
