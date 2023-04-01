package database_3ds

import (
	"database/sql"

	"github.com/PretendoNetwork/friends-secure/database"
	"github.com/PretendoNetwork/friends-secure/globals"
	"github.com/PretendoNetwork/nex-go"
)

// Update a user's last online time
func UpdateUserLastOnlineTime(pid uint32, lastOnline *nex.DateTime) {
	var showOnline bool

	err := database.Postgres.QueryRow(`SELECT show_online FROM "3ds".user_data WHERE pid=$1`, pid).Scan(&showOnline)
	if err != nil && err != sql.ErrNoRows {
		globals.Logger.Critical(err.Error())
	}
	if !showOnline {
		return
	}

	_, err = database.Postgres.Exec(`
		INSERT INTO "3ds".user_data (pid, last_online)
		VALUES ($1, $2)
		ON CONFLICT (pid)
		DO UPDATE SET 
		last_online = $2`, pid, lastOnline.Value())
	if err != nil {
		globals.Logger.Critical(err.Error())
	}
}
