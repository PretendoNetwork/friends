package database_3ds

import (
	"github.com/PretendoNetwork/friends/database"
)

// UpdateUserPreferences updates a user's preferences
func UpdateUserPreferences(pid uint32, show_online bool, show_current_game bool) error {
	_, err := database.Postgres.Exec(`
		INSERT INTO "3ds".user_data (pid, show_online, show_current_game)
		VALUES ($1, $2, $3)
		ON CONFLICT (pid)
		DO UPDATE SET 
		show_online = $2,
		show_current_game = $3`, pid, show_online, show_current_game)

	if err != nil {
		return err
	}

	return nil
}
