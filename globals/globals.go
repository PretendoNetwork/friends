package globals

import (
	"github.com/PretendoNetwork/friends/types"
	pb "github.com/PretendoNetwork/grpc-go/account"
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/plogger-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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
var AESKey []byte
var GRPCAccountClientConnection *grpc.ClientConn
var GRPCAccountClient pb.AccountClient
var GRPCAccountCommonMetadata metadata.MD
var DatabaseMaxConnections int
