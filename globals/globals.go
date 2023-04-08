package globals

import (
	"crypto/rsa"

	"github.com/PretendoNetwork/friends-secure/types"
	"github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/plogger-go"
)

var Logger = plogger.NewLogger()
var NEXServer *nex.Server
var ConnectedUsers map[uint32]*types.ConnectedUser
var RSAPrivateKeyBytes []byte
var RSAPrivateKey *rsa.PrivateKey
var HMACSecret []byte
