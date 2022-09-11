package database

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"

	"github.com/PretendoNetwork/friends-secure/globals"
)

var Postgres *sql.DB

func connectPostgres() {
	var err error

	Postgres, err = sql.Open("postgres", os.Getenv("DATABASE_URI"))
	if err != nil {
		globals.Logger.Critical(err.Error())
	}

	globals.Logger.Success("Connected to Postgres!")

	initPostgresWiiU()
	// TODO: 3DS database
}
