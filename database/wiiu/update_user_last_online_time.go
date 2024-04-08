package database_wiiu

import (
	"github.com/PretendoNetwork/friends/database"
	"github.com/PretendoNetwork/nex-go/v2/types"
)

// UpdateUserLastOnlineTime updates a user's last online time
func UpdateUserLastOnlineTime(pid uint32, lastOnline *types.DateTime) error {
	_, err := database.Postgres.Exec(`
		INSERT INTO wiiu.user_data (pid, last_online)
		VALUES ($1, $2)
		ON CONFLICT (pid)
		DO UPDATE SET 
		last_online = $2`, pid, lastOnline.Value())

	if err != nil {
		return err
	}

	return nil
}
