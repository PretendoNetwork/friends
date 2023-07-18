package nex

import (
	"github.com/PretendoNetwork/friends-secure/globals"
	secureconnection "github.com/PretendoNetwork/nex-protocols-common-go/secure-connection"
)

func registerCommonSecureServerProtocols() {
	secureconnection.NewCommonSecureConnectionProtocol(globals.SecureServer)
}
