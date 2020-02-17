package main

import (
	//"fmt"

	"fmt"

	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

var nexServer *nex.Server
var secureServer *nexproto.SecureProtocol

func main() {
	nexServer = nex.NewServer()
	nexServer.SetPrudpVersion(0)
	nexServer.SetSignatureVersion(1)
	nexServer.SetKerberosKeySize(16)
	nexServer.SetAccessKey("ridfebb9")

	secureServer = nexproto.NewSecureProtocol(nexServer)
	accountManagementServer := nexproto.NewAccountManagementProtocol(nexServer)
	friendsServer := nexproto.NewFriendsProtocol(nexServer)

	// Handle PRUDP CONNECT packet (not an RMC method)
	nexServer.On("Connect", connect)

	// Account Management protocol handles

	accountManagementServer.NintendoCreateAccount(nintendoCreateAccount)

	// Secure protocol handles

	// Handle RegisterEx RMC method
	secureServer.RegisterEx(registerEx)

	// Friends (WiiU) protocol handles

	friendsServer.UpdateAndGetAllInformation(updateAndGetAllInformation)

	friendsServer.CheckSettingStatus(checkSettingStatus)

	nexServer.Listen(":60001")
}
