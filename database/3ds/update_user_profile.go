package database_3ds

import (
	"github.com/PretendoNetwork/friends/database"
	friends_3ds_types "github.com/PretendoNetwork/nex-protocols-go/v2/friends-3ds/types"
)

// UpdateUserProfile updates a user's profile
func UpdateUserProfile(pid uint32, profileData friends_3ds_types.MyProfile) error {
	_, err := database.Manager.Exec(`
		INSERT INTO "3ds".user_data (pid, region, area, language, country)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (pid)
		DO UPDATE SET 
		region = $2,
		area = $3,
		language = $4,
		country = $5`, pid, uint8(profileData.Region), uint8(profileData.Area), uint8(profileData.Language), uint8(profileData.Country))

	if err != nil {
		return err
	}

	return nil
}
