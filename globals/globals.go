package globals

import (
	"github.com/PretendoNetwork/friends/types"
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/plogger-go"
	"google.golang.org/grpc"
)

var Logger *plogger.Logger
var AuthenticationServerAccount *nex.Account
var SecureServerAccount *nex.Account
var GuestAccount *nex.Account
var KerberosPassword = "password" // * Default password
var AuthenticationServer *nex.PRUDPServer
var AuthenticationEndpoint *nex.PRUDPEndPoint
var SecureServer *nex.PRUDPServer
var SecureEndpoint *nex.PRUDPEndPoint
var ConnectedUsers *nex.MutexMap[uint32, *types.ConnectedUser]
var GRPCAccountClientConnection *grpc.ClientConn
