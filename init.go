package main

import (
	"cmp"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/PretendoNetwork/friends/database"
	"github.com/PretendoNetwork/friends/globals"
	"github.com/PretendoNetwork/friends/types"
	pb "github.com/PretendoNetwork/grpc-go/account"
	"github.com/PretendoNetwork/nex-go/v2"
	nex_types "github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/PretendoNetwork/plogger-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/joho/godotenv"
)

func init() {
	globals.Logger = plogger.NewLogger()
	globals.ConnectedUsers = nex.NewMutexMap[uint32, *types.ConnectedUser]()

	// * Setup RSA private key for token parsing
	var err error

	err = godotenv.Load()
	if err != nil {
		globals.Logger.Warningf("Error loading .env file: %s", err.Error())
	}

	postgresURI := os.Getenv("PN_FRIENDS_CONFIG_DATABASE_URI")
	databaseMaxConnectionsStr := cmp.Or(os.Getenv("PN_FRIENDS_CONFIG_DATABASE_MAX_CONNECTIONS"), "100")
	aesKey := os.Getenv("PN_FRIENDS_CONFIG_AES_KEY")
	grpcAPIKey := os.Getenv("PN_FRIENDS_CONFIG_GRPC_API_KEY")
	grpcServerPort := os.Getenv("PN_FRIENDS_GRPC_SERVER_PORT")
	authenticationServerPort := os.Getenv("PN_FRIENDS_AUTHENTICATION_SERVER_PORT")
	secureServerHost := os.Getenv("PN_FRIENDS_SECURE_SERVER_HOST")
	secureServerPort := os.Getenv("PN_FRIENDS_SECURE_SERVER_PORT")
	accountGRPCHost := os.Getenv("PN_FRIENDS_ACCOUNT_GRPC_HOST")
	accountGRPCPort := os.Getenv("PN_FRIENDS_ACCOUNT_GRPC_PORT")
	accountGRPCAPIKey := os.Getenv("PN_FRIENDS_ACCOUNT_GRPC_API_KEY")

	if strings.TrimSpace(postgresURI) == "" {
		globals.Logger.Error("PN_FRIENDS_CONFIG_DATABASE_URI environment variable not set")
		os.Exit(0)
	}

	databaseMaxConnections, err := strconv.Atoi(databaseMaxConnectionsStr)

	if err != nil {
		globals.Logger.Errorf("PN_FRIENDS_CONFIG_DATABASE_MAX_CONNECTIONS is not a valid number. Got %s", databaseMaxConnectionsStr)
		os.Exit(0)
	} else {
		globals.DatabaseMaxConnections = databaseMaxConnections
	}

	kerberosPassword := make([]byte, 0x10)
	_, err = rand.Read(kerberosPassword)
	if err != nil {
		globals.Logger.Error("Error generating Kerberos password")
		os.Exit(0)
	}

	globals.KerberosPassword = string(kerberosPassword)

	globals.AuthenticationServerAccount = nex.NewAccount(nex_types.NewPID(1), "Quazal Authentication", globals.KerberosPassword)
	globals.SecureServerAccount = nex.NewAccount(nex_types.NewPID(2), "Quazal Rendez-Vous", globals.KerberosPassword)
	globals.GuestAccount = nex.NewAccount(nex_types.NewPID(100), "guest", "MMQea3n!fsik") // * Guest account password is always the same, known to all consoles. Only allow on the friends server

	if strings.TrimSpace(aesKey) == "" {
		globals.Logger.Error("PN_FRIENDS_CONFIG_AES_KEY environment variable not set")
		os.Exit(0)
	} else {
		globals.AESKey, err = hex.DecodeString(os.Getenv("PN_FRIENDS_CONFIG_AES_KEY"))
		if err != nil {
			globals.Logger.Criticalf("Failed to decode AES key: %v", err)
			os.Exit(0)
		}
	}

	if strings.TrimSpace(grpcAPIKey) == "" {
		globals.Logger.Warning("Insecure gRPC server detected. PN_FRIENDS_CONFIG_GRPC_API_KEY environment variable not set")
	}

	if strings.TrimSpace(grpcServerPort) == "" {
		globals.Logger.Error("PN_FRIENDS_GRPC_SERVER_PORT environment variable not set")
		os.Exit(0)
	}

	if port, err := strconv.Atoi(grpcServerPort); err != nil {
		globals.Logger.Errorf("PN_FRIENDS_GRPC_SERVER_PORT is not a valid port. Expected 0-65535, got %s", grpcServerPort)
		os.Exit(0)
	} else if port < 0 || port > 65535 {
		globals.Logger.Errorf("PN_FRIENDS_GRPC_SERVER_PORT is not a valid port. Expected 0-65535, got %s", grpcServerPort)
		os.Exit(0)
	}

	if strings.TrimSpace(authenticationServerPort) == "" {
		globals.Logger.Error("PN_FRIENDS_AUTHENTICATION_SERVER_PORT environment variable not set")
		os.Exit(0)
	}

	if port, err := strconv.Atoi(authenticationServerPort); err != nil {
		globals.Logger.Errorf("PN_FRIENDS_AUTHENTICATION_SERVER_PORT is not a valid port. Expected 0-65535, got %s", authenticationServerPort)
		os.Exit(0)
	} else if port < 0 || port > 65535 {
		globals.Logger.Errorf("PN_FRIENDS_AUTHENTICATION_SERVER_PORT is not a valid port. Expected 0-65535, got %s", authenticationServerPort)
		os.Exit(0)
	}

	if strings.TrimSpace(secureServerHost) == "" {
		globals.Logger.Error("PN_FRIENDS_SECURE_SERVER_HOST environment variable not set")
		os.Exit(0)
	}

	if strings.TrimSpace(secureServerPort) == "" {
		globals.Logger.Error("PN_FRIENDS_SECURE_SERVER_PORT environment variable not set")
		os.Exit(0)
	}

	if port, err := strconv.Atoi(secureServerPort); err != nil {
		globals.Logger.Errorf("PN_FRIENDS_SECURE_SERVER_PORT is not a valid port. Expected 0-65535, got %s", secureServerPort)
		os.Exit(0)
	} else if port < 0 || port > 65535 {
		globals.Logger.Errorf("PN_FRIENDS_SECURE_SERVER_PORT is not a valid port. Expected 0-65535, got %s", secureServerPort)
		os.Exit(0)
	}

	if strings.TrimSpace(accountGRPCHost) == "" {
		globals.Logger.Error("PN_FRIENDS_ACCOUNT_GRPC_HOST environment variable not set")
		os.Exit(0)
	}

	if strings.TrimSpace(accountGRPCPort) == "" {
		globals.Logger.Error("PN_FRIENDS_ACCOUNT_GRPC_PORT environment variable not set")
		os.Exit(0)
	}

	if port, err := strconv.Atoi(accountGRPCPort); err != nil {
		globals.Logger.Errorf("PN_FRIENDS_ACCOUNT_GRPC_PORT is not a valid port. Expected 0-65535, got %s", accountGRPCPort)
		os.Exit(0)
	} else if port < 0 || port > 65535 {
		globals.Logger.Errorf("PN_FRIENDS_ACCOUNT_GRPC_PORT is not a valid port. Expected 0-65535, got %s", accountGRPCPort)
		os.Exit(0)
	}

	if strings.TrimSpace(accountGRPCAPIKey) == "" {
		globals.Logger.Warning("Insecure gRPC server detected. PN_FRIENDS_ACCOUNT_GRPC_API_KEY environment variable not set")
	}

	globals.GRPCAccountClientConnection, err = grpc.Dial(fmt.Sprintf("%s:%s", accountGRPCHost, accountGRPCPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		globals.Logger.Criticalf("Failed to connect to account gRPC server: %v", err)
		os.Exit(0)
	}

	globals.GRPCAccountClient = pb.NewAccountClient(globals.GRPCAccountClientConnection)
	globals.GRPCAccountCommonMetadata = metadata.Pairs(
		"X-API-Key", accountGRPCAPIKey,
	)

	database.ConnectPostgres()
}
