package database_3ds

import (
	"github.com/PretendoNetwork/friends-secure/database"
	"github.com/PretendoNetwork/friends-secure/globals"
	friends_3ds "github.com/PretendoNetwork/nex-protocols-go/friends/3ds"
)

// Update a user's profile
func UpdateUserProfile(pid uint32, profileData *friends_3ds.MyProfile) {
	_, err := database.Postgres.Exec(`
		INSERT INTO "3ds".user_data (pid, region, area, language)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (pid)
		DO UPDATE SET 
		region = $2,
		area = $3,
		language = $4`, pid, profileData.Region, profileData.Area, profileData.Language)

	if err != nil {
		globals.Logger.Critical(err.Error())
	}
}
