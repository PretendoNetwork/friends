package main

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"os"
	"strings"

	"github.com/PretendoNetwork/friends/database"
	"github.com/PretendoNetwork/friends/globals"
	"github.com/PretendoNetwork/friends/types"
	"github.com/PretendoNetwork/nex-go/v2"
	nex_types "github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	"github.com/PretendoNetwork/plogger-go"

	"github.com/joho/godotenv"
)

func init() {
	globals.Logger = plogger.NewLogger()
	globals.ConnectedUsers = nex.NewMutexMap[uint32, *types.ConnectedUser]()

	var err error

	err = godotenv.Load()
	if err != nil {
		globals.Logger.Warningf("Error loading .env file: %s", err.Error())
	}

	globals.Config = globals.NewConfigParser(globals.Config).SetPrefix("PN_FRIENDS_CONFIG").ParseFromEnv()

	kerberosPassword := make([]byte, 0x10)
	_, err = rand.Read(kerberosPassword)
	if err != nil {
		globals.Logger.Error("Error generating Kerberos password")
		os.Exit(0)
	}

	globals.KerberosPassword = string(kerberosPassword)

	globals.AuthenticationServerAccount = nex.NewAccount(nex_types.NewPID(1), "Quazal Authentication", globals.KerberosPassword, false)
	globals.SecureServerAccount = nex.NewAccount(nex_types.NewPID(2), "Quazal Rendez-Vous", globals.KerberosPassword, false)
	globals.GuestAccount = nex.NewAccount(nex_types.NewPID(100), "guest", "MMQea3n!fsik", false)

	if strings.TrimSpace(globals.Config.GRPCAPIKey) == "" {
		globals.Logger.Warning("Insecure gRPC server detected. PN_FRIENDS_CONFIG_GRPC_API_KEY environment variable not set")
	}

	if strings.TrimSpace(globals.Config.AccountGRPCAPIKey) == "" {
		globals.Logger.Warning("Insecure gRPC server detected. PN_FRIENDS_CONFIG_ACCOUNT_GRPC_API_KEY environment variable not set")
	}

	common_globals.ConnectToAccountGRPC(globals.Config.AccountGRPCHost, globals.Config.AccountGRPCPort, globals.Config.AccountGRPCAPIKey)

	if strings.TrimSpace(globals.Config.MiiDecryptKey) == "" {
		globals.Logger.Warning("PN_FRIENDS_CONFIG_MII_DECRYPT_KEY environment variable not set. 3DS Mii data cannot be decrypted")
	}

	miiKeyBytes, err := hex.DecodeString(globals.Config.MiiDecryptKey)
	if err != nil {
		globals.Logger.Criticalf("Failed to decode PN_FRIENDS_CONFIG_MII_DECRYPT_KEY %v", err)
		os.Exit(0)
	}

	miiMD5Hash := md5.Sum(miiKeyBytes)
	if hex.EncodeToString(miiMD5Hash[:]) != "aeb707b225ec0fcd8a503e26e3dcd596" {
		globals.Logger.Criticalf("PN_FRIENDS_CONFIG_MII_DECRYPT_KEY is incorrect! md5: %s", hex.EncodeToString(miiMD5Hash[:]))
		os.Exit(0)
	}

	database.ConnectPostgres()
}
