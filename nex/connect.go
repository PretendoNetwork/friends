package nex

import (
	"time"

	"github.com/PretendoNetwork/friends-secure/globals"
	"github.com/PretendoNetwork/friends-secure/types"
	nex "github.com/PretendoNetwork/nex-go"
)

func connect(packet *nex.PacketV0) {
	// * We aren't making any replies here because the common Secure Protocol already does that
	// *
	// * We only want to check that the data given is right so that we don't register a client
	// * with an invalid request
	payload := packet.Payload()
	stream := nex.NewStreamIn(payload, globals.SecureServer)

	ticketData, err := stream.ReadBuffer()
	if err != nil {
		return
	}

	requestData, err := stream.ReadBuffer()
	if err != nil {
		return
	}

	serverKey := nex.DeriveKerberosKey(2, []byte(globals.SecureServer.KerberosPassword()))

	ticket := nex.NewKerberosTicketInternalData()
	err = ticket.Decrypt(nex.NewStreamIn(ticketData, globals.SecureServer), serverKey)
	if err != nil {
		return
	}

	ticketTime := ticket.Timestamp().Standard()
	serverTime := time.Now().UTC()

	timeLimit := ticketTime.Add(time.Minute * 2)
	if serverTime.After(timeLimit) {
		return
	}

	sessionKey := ticket.SessionKey()
	kerberos, err := nex.NewKerberosEncryption(sessionKey)
	if err != nil {
		return
	}

	decryptedRequestData := kerberos.Decrypt(requestData)
	checkDataStream := nex.NewStreamIn(decryptedRequestData, globals.SecureServer)

	userPID, err := checkDataStream.ReadUInt32LE()
	if err != nil {
		return
	}

	_, err = checkDataStream.ReadUInt32LE() // CID of secure server station url
	if err != nil {
		return
	}

	_, err = checkDataStream.ReadUInt32LE() // Response check
	if err != nil {
		return
	}

	connectedUser := types.NewConnectedUser()
	connectedUser.PID = userPID
	connectedUser.Client = packet.Sender()
	globals.ConnectedUsers[userPID] = connectedUser
}
