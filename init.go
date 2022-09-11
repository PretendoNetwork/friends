package main

import (
	"crypto/rsa"
	"io/ioutil"
	"log"

	"github.com/PretendoNetwork/friends-secure/database"
	"github.com/PretendoNetwork/friends-secure/globals"

	"github.com/joho/godotenv"
)

/*
type Config struct {
	Mongo struct {
	}
	Cassandra struct{}
}
*/

type nexToken struct {
	SystemType  uint8
	TokenType   uint8
	UserPID     uint32
	AccessLevel uint8
	TitleID     uint64
	ExpireTime  uint64
}

var rsaPrivateKeyBytes []byte
var rsaPrivateKey *rsa.PrivateKey
var hmacSecret []byte

func init() {
	globals.ConnectedUsers = make(map[uint32]*globals.ConnectedUser)
	// Setup RSA private key for token parsing
	var err error

	rsaPrivateKeyBytes, err = ioutil.ReadFile("private.pem")
	if err != nil {
		// TODO: Handle error
		globals.Logger.Critical(err.Error())
	}

	rsaPrivateKey, err = parseRsaPrivateKey(rsaPrivateKeyBytes)
	if err != nil {
		// TODO: Handle error
		globals.Logger.Critical(err.Error())
	}

	hmacSecret, err = ioutil.ReadFile("secret.key")
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
