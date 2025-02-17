package nex

import (
	"os"
	"strconv"

	"github.com/PretendoNetwork/friends/globals"
	"github.com/PretendoNetwork/nex-go/v2/constants"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_ticket_granting "github.com/PretendoNetwork/nex-protocols-common-go/v2/ticket-granting"
	ticket_granting "github.com/PretendoNetwork/nex-protocols-go/v2/ticket-granting"
)

func registerCommonAuthenticationServerProtocols() {
	ticketGrantingProtocol := ticket_granting.NewProtocol()
	commonTicketGrantingProtocol := common_ticket_granting.NewCommonProtocol(ticketGrantingProtocol)

	port, _ := strconv.Atoi(os.Getenv("PN_FRIENDS_SECURE_SERVER_PORT"))

	secureStationURL := types.NewStationURL("")
	secureStationURL.SetURLType(constants.StationURLPRUDPS)
	secureStationURL.SetAddress(os.Getenv("PN_FRIENDS_SECURE_SERVER_HOST"))
	secureStationURL.SetPortNumber(uint16(port))
	secureStationURL.SetConnectionID(1)
	secureStationURL.SetPrincipalID(types.NewPID(2))
	secureStationURL.SetStreamID(1)
	secureStationURL.SetStreamType(constants.StreamTypeRVSecure)
	secureStationURL.SetType(2)

	commonTicketGrantingProtocol.SecureServerAccount = globals.SecureEndpoint.ServerAccount
	commonTicketGrantingProtocol.SessionKeyLength = 16
	commonTicketGrantingProtocol.SecureStationURL = secureStationURL
	commonTicketGrantingProtocol.BuildName = types.NewString(serverBuildString)
	commonTicketGrantingProtocol.EnableInsecureLogin()

	globals.AuthenticationEndpoint.RegisterServiceProtocol(ticketGrantingProtocol)
}
