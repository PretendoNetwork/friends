package nex

import (
	"fmt"
	"os"
	"strconv"

	"github.com/PretendoNetwork/friends/globals"
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
)

var serverBuildString string

func StartAuthenticationServer() {
	port, _ := strconv.Atoi(os.Getenv("PN_FRIENDS_AUTHENTICATION_SERVER_PORT"))

	globals.AuthenticationServer = nex.NewPRUDPServer()
	globals.AuthenticationEndpoint = nex.NewPRUDPEndPoint(1)

	globals.AuthenticationEndpoint.AccountDetailsByPID = globals.AccountDetailsByPID
	globals.AuthenticationEndpoint.AccountDetailsByUsername = globals.AccountDetailsByUsername
	globals.AuthenticationEndpoint.ServerAccount = nex.NewAccount(types.NewPID(1), "Quazal Authentication", os.Getenv("PN_FRIENDS_CONFIG_AUTHENTICATION_PASSWORD"))

	globals.AuthenticationEndpoint.OnData(func(packet nex.PacketInterface) {
		request := packet.RMCMessage()

		fmt.Println("==Friends - Auth==")
		fmt.Printf("Protocol ID: %#v\n", request.ProtocolID)
		fmt.Printf("Method ID: %#v\n", request.MethodID)
		fmt.Println("===============")
	})

	registerCommonAuthenticationServerProtocols()

	globals.AuthenticationServer.SetFragmentSize(962)
	globals.AuthenticationServer.LibraryVersions.SetDefault(nex.NewLibraryVersion(1, 1, 0))
	globals.AuthenticationServer.SessionKeyLength = 16
	globals.AuthenticationServer.AccessKey = "ridfebb9"
	globals.AuthenticationServer.BindPRUDPEndPoint(globals.AuthenticationEndpoint)
	globals.AuthenticationServer.Listen(port)
}
