package main

import (
	"crypto/rsa"
	"io/ioutil"
)

type Config struct {
	Mongo struct {
	}
	Cassandra struct{}
}

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

func init() {
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

	//connectMongo()
}
