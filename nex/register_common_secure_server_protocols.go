package nex

import (
	"github.com/PretendoNetwork/friends/globals"
	nex_secure_connection "github.com/PretendoNetwork/friends/nex/secure-connection"
	secure_connection "github.com/PretendoNetwork/nex-protocols-go/secure-connection"
	common_secure_connection "github.com/PretendoNetwork/nex-protocols-common-go/secure-connection"
)

func registerCommonSecureServerProtocols() {
	secureConnectionProtocol := secure_connection.NewProtocol(globals.SecureServer)
	common_secure_connection.NewCommonSecureConnectionProtocol(secureConnectionProtocol)

	secureConnectionProtocol.RegisterEx = nex_secure_connection.RegisterEx
}
