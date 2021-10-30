package main

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"log"
	"runtime"

	"github.com/bwmarrin/snowflake"
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
	SystemType uint8
	TokenType  uint8
	UserPID    uint32
	TitleID    uint64
	CreatTime  uint64
}

var rsaPrivateKeyBytes []byte
var rsaPrivateKey *rsa.PrivateKey
var hmacSecret []byte
var snowflakeNodes []*snowflake.Node

func init() {
	// Setup RSA private key for token parsing
	var err error

	rsaPrivateKeyBytes, err = ioutil.ReadFile("private.pem")
	if err != nil {
		panic(err)
	}

	rsaPrivateKey, err = parseRsaPrivateKey(rsaPrivateKeyBytes)
	if err != nil {
		panic(err)
	}

	hmacSecret, err = ioutil.ReadFile("secret.key")
	if err != nil {
		panic(err)
	}

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connectMongo()
	connectCassandra()
	createSnowflakeNodes()
}

func createSnowflakeNodes() {
	snowflakeNodes = make([]*snowflake.Node, 0)

	for corenum := 0; corenum < runtime.NumCPU(); corenum++ {
		node, err := snowflake.NewNode(int64(corenum))
		if err != nil {
			fmt.Println(err)
			return
		}
		snowflakeNodes = append(snowflakeNodes, node)
	}
}
