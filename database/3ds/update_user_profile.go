package database_3ds

import (
	"github.com/PretendoNetwork/friends/database"
	friends_3ds_types "github.com/PretendoNetwork/nex-protocols-go/friends-3ds/types"
)

// UpdateUserProfile updates a user's profile
func UpdateUserProfile(pid uint32, profileData *friends_3ds_types.MyProfile) error {
	_, err := database.Postgres.Exec(`
		INSERT INTO "3ds".user_data (pid, region, area, language)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (pid)
		DO UPDATE SET 
		region = $2,
		area = $3,
		language = $4`, pid, profileData.Region, profileData.Area, profileData.Language)

	if err != nil {
		return err
	}

	return nil
}
