package nex

import (
	"os"
	"strconv"

	"github.com/PretendoNetwork/friends/globals"
	nex "github.com/PretendoNetwork/nex-go"
	ticket_granting "github.com/PretendoNetwork/nex-protocols-common-go/ticket-granting"
)

func registerCommonAuthenticationServerProtocols() {
	ticketGrantingProtocol := ticket_granting.NewCommonTicketGrantingProtocol(globals.AuthenticationServer)

	port, _ := strconv.Atoi(os.Getenv("PN_FRIENDS_SECURE_SERVER_PORT"))

	secureStationURL := nex.NewStationURL("")
	secureStationURL.SetScheme("prudps")
	secureStationURL.SetAddress(os.Getenv("PN_FRIENDS_SECURE_SERVER_HOST"))
	secureStationURL.SetPort(uint32(port))
	secureStationURL.SetCID(1)
	secureStationURL.SetPID(nex.NewPID[uint32](2))
	secureStationURL.SetSID(1)
	secureStationURL.SetStream(10)
	secureStationURL.SetType(2)

	ticketGrantingProtocol.SetSecureStationURL(secureStationURL)
	ticketGrantingProtocol.SetBuildName(serverBuildString)
	ticketGrantingProtocol.EnableInsecureLogin()

	globals.AuthenticationServer.PasswordFromPID = globals.PasswordFromPID
}
