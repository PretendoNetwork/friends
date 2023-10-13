package nex

import (
	"github.com/PretendoNetwork/friends/globals"
	nex_secure_connection "github.com/PretendoNetwork/friends/nex/secure-connection"
	secureconnection "github.com/PretendoNetwork/nex-protocols-common-go/secure-connection"
)

func registerCommonSecureServerProtocols() {
	secureConnectionProtocol := secureconnection.NewCommonSecureConnectionProtocol(globals.SecureServer)
	secureConnectionProtocol.RegisterEx(nex_secure_connection.RegisterEx)
}
