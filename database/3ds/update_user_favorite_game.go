package database_3ds

import (
	"github.com/PretendoNetwork/friends-secure/database"
	friends_3ds_types "github.com/PretendoNetwork/nex-protocols-go/friends-3ds/types"
)

// UpdateUserFavoriteGame updates a user's favorite game
func UpdateUserFavoriteGame(pid uint32, gameKey *friends_3ds_types.GameKey) error {
	_, err := database.Postgres.Exec(`
		INSERT INTO "3ds".user_data (pid, favorite_title, favorite_title_version)
		VALUES ($1, $2, $3)
		ON CONFLICT (pid)
		DO UPDATE SET 
		favorite_title = $2,
		favorite_title_version = $3`, pid, gameKey.TitleID, gameKey.TitleVersion)

	if err != nil {
		return err
	}

	return nil
}
