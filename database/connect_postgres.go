package database

import (
	"os"

	_ "github.com/lib/pq"

	"github.com/PretendoNetwork/friends/globals"
	sqlmanager "github.com/PretendoNetwork/sql-manager"
)

var Manager *sqlmanager.SQLManager

func ConnectPostgres() {
	var err error

	Manager, err = sqlmanager.NewSQLManager("postgres", os.Getenv("PN_FRIENDS_CONFIG_DATABASE_URI"), int64(globals.DatabaseMaxConnections))
	if err != nil {
		globals.Logger.Critical(err.Error())
	}

	globals.Logger.Success("Connected to Postgres!")

	initPostgresWiiU()
	initPostgres3DS()
}
