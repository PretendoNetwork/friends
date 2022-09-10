package database

import (
	"github.com/PretendoNetwork/plogger-go"
)

var logger = plogger.NewLogger()

func Connect() {
	connectMongo()
	connectCassandra()
}
