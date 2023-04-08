package database_wiiu

import (
	"database/sql"

	"github.com/PretendoNetwork/friends-secure/database"
	"github.com/PretendoNetwork/friends-secure/globals"
	friends_wiiu "github.com/PretendoNetwork/nex-protocols-go/friends/wiiu"
)

func GetUserPrincipalPreference(pid uint32) *friends_wiiu.PrincipalPreference {
	preference := friends_wiiu.NewPrincipalPreference()

	err := database.Postgres.QueryRow(`SELECT show_online, show_current_game, block_friend_requests FROM wiiu.user_data WHERE pid=$1`, pid).Scan(&preference.ShowOnlinePresence, &preference.ShowCurrentTitle, &preference.BlockFriendRequests)
	if err != nil {
		if err == sql.ErrNoRows {
			globals.Logger.Warning(err.Error())
		} else {
			globals.Logger.Critical(err.Error())
		}
	}

	return preference
}
