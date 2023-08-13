package nex

import (
	"fmt"
	"os"

	"github.com/PretendoNetwork/friends/globals"
	"github.com/PretendoNetwork/nex-go"
)

var serverBuildString string

func StartAuthenticationServer() {
	globals.AuthenticationServer = nex.NewServer()
	globals.AuthenticationServer.SetPRUDPVersion(0)
	globals.AuthenticationServer.SetPRUDPProtocolMinorVersion(0) // TODO: Figure out what to put here
	globals.AuthenticationServer.SetDefaultNEXVersion(&nex.NEXVersion{
		Major: 1,
		Minor: 1,
		Patch: 0,
	})
	globals.AuthenticationServer.SetKerberosKeySize(16)
	globals.AuthenticationServer.SetKerberosPassword(globals.KerberosPassword)
	globals.AuthenticationServer.SetAccessKey("ridfebb9")

	globals.AuthenticationServer.On("Data", func(packet *nex.PacketV0) {
		request := packet.RMCRequest()

		fmt.Println("==Friends - Auth==")
		fmt.Printf("Protocol ID: %#v\n", request.ProtocolID())
		fmt.Printf("Method ID: %#v\n", request.MethodID())
		fmt.Println("===============")
	})

	registerCommonAuthenticationServerProtocols()

	globals.AuthenticationServer.Listen(fmt.Sprintf(":%s", os.Getenv("PN_FRIENDS_AUTHENTICATION_SERVER_PORT")))
}
