package main

import (
	"encoding/hex"
	"log"
	"os"

	"github.com/PretendoNetwork/friends-secure/database"
	"github.com/PretendoNetwork/friends-secure/globals"
	"github.com/PretendoNetwork/friends-secure/types"

	"github.com/joho/godotenv"
)

func init() {
	globals.ConnectedUsers = make(map[uint32]*types.ConnectedUser)
	// Setup RSA private key for token parsing
	var err error

	err = godotenv.Load()
	if err != nil {
		globals.Logger.Warning("Error loading .env file")
	}

	globals.AESKey, err = hex.DecodeString(os.Getenv("PN_FRIENDS_CONFIG_AES_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	database.Connect()
}
