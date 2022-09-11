package database_wiiu

import (
	"github.com/PretendoNetwork/friends-secure/database"
	"github.com/PretendoNetwork/friends-secure/globals"
	"github.com/PretendoNetwork/nex-go"
)

func UpdateUserLastOnlineTime(pid uint32, lastOnline *nex.DateTime) {
	_, err := database.Postgres.Exec(`
		INSERT INTO wiiu.user_data (pid, last_online)
		VALUES ($1, $2)
		ON CONFLICT (pid)
		DO UPDATE SET 
		last_online = $2`, pid, lastOnline.Value())

	if err != nil {
		globals.Logger.Critical(err.Error())
	}
}
