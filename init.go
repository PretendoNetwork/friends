package main

import (
	"io/ioutil"
	"log"

	"github.com/PretendoNetwork/friends-secure/database"
	"github.com/PretendoNetwork/friends-secure/globals"
	"github.com/PretendoNetwork/friends-secure/types"
	"github.com/PretendoNetwork/friends-secure/utility"

	"github.com/joho/godotenv"
)

func init() {
	globals.ConnectedUsers = make(map[uint32]*types.ConnectedUser)
	// Setup RSA private key for token parsing
	var err error

	globals.RSAPrivateKeyBytes, err = ioutil.ReadFile("private.pem")
	if err != nil {
		// TODO: Handle error
		globals.Logger.Critical(err.Error())
	}

	globals.RSAPrivateKey, err = utility.ParseRsaPrivateKey(globals.RSAPrivateKeyBytes)
	if err != nil {
		// TODO: Handle error
		globals.Logger.Critical(err.Error())
	}

	globals.HMACSecret, err = ioutil.ReadFile("secret.key")
	if err != nil {
		// TODO: Handle error
		globals.Logger.Critical(err.Error())
	}

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.Connect()
}
