package database_3ds

import (
	"database/sql"

	"github.com/PretendoNetwork/friends/database"
	"github.com/PretendoNetwork/nex-go/v2/types"
)

// UpdateUserLastOnlineTime updates a user's last online time
func UpdateUserLastOnlineTime(pid uint32, lastOnline *types.DateTime) error {
	var showOnline bool

	err := database.Postgres.QueryRow(`SELECT show_online FROM "3ds".user_data WHERE pid=$1`, pid).Scan(&showOnline)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if !showOnline {
		return nil
	}

	_, err = database.Postgres.Exec(`
		INSERT INTO "3ds".user_data (pid, last_online)
		VALUES ($1, $2)
		ON CONFLICT (pid)
		DO UPDATE SET 
		last_online = $2`, pid, lastOnline.Value())
	if err != nil {
		return err
	}

	return nil
}
