package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/PretendoNetwork/friends-secure/database"
	"github.com/PretendoNetwork/friends-secure/globals"
	"github.com/PretendoNetwork/friends-secure/types"
	pb "github.com/PretendoNetwork/grpc-go/account"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/joho/godotenv"
)

func init() {
	globals.ConnectedUsers = make(map[uint32]*types.ConnectedUser)
	// Setup RSA private key for token parsing
	var err error

	err = godotenv.Load()
	if err != nil {
		globals.Logger.Warningf("Error loading .env file: %s", err.Error())
	}

	globals.AESKey, err = hex.DecodeString(os.Getenv("PN_FRIENDS_CONFIG_AES_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	globals.GRPCAccountClientConnection, err = grpc.Dial(fmt.Sprintf("%s:%s", os.Getenv("PN_FRIENDS_ACCOUNT_GRPC_HOST"), os.Getenv("PN_FRIENDS_ACCOUNT_GRPC_PORT")), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		globals.Logger.Criticalf("Failed to connect to account gRPC server: %v", err)
		os.Exit(0)
	}

	globals.GRPCAccountClient = pb.NewAccountClient(globals.GRPCAccountClientConnection)
	globals.GRPCAccountCommonMetadata = metadata.Pairs(
		"X-API-Key", os.Getenv("PN_FRIENDS_ACCOUNT_GRPC_APIKEY"),
	)

	globals.KerberosPassword = os.Getenv("PN_FRIENDS_CONFIG_KERBEROS_PASSWORD")

	database.Connect()
}
