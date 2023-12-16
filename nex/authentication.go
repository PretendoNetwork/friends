package nex

import (
	"fmt"
	"os"
	"strconv"

	"github.com/PretendoNetwork/friends/globals"
	"github.com/PretendoNetwork/nex-go"
)

var serverBuildString string

func StartAuthenticationServer() {
	globals.AuthenticationServer = nex.NewPRUDPServer()
	globals.AuthenticationServer.SetFragmentSize(962)
	globals.AuthenticationServer.SetDefaultLibraryVersion(nex.NewLibraryVersion(1, 1, 0))
	globals.AuthenticationServer.SetKerberosPassword([]byte(globals.KerberosPassword))
	globals.AuthenticationServer.SetKerberosKeySize(16)
	globals.AuthenticationServer.SetAccessKey("ridfebb9")

	globals.AuthenticationServer.OnData(func(packet nex.PacketInterface) {
		request := packet.RMCMessage()

		fmt.Println("==Friends - Auth==")
		fmt.Printf("Protocol ID: %#v\n", request.ProtocolID)
		fmt.Printf("Method ID: %#v\n", request.MethodID)
		fmt.Println("===============")
	})

	registerCommonAuthenticationServerProtocols()

	port, _ := strconv.Atoi(os.Getenv("PN_FRIENDS_AUTHENTICATION_SERVER_PORT"))
	globals.AuthenticationServer.Listen(port)
}
