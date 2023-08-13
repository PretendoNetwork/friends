package database_wiiu

import (
	"database/sql"

	"github.com/PretendoNetwork/friends/database"
	"github.com/PretendoNetwork/friends/globals"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/friends-wiiu/types"
)

func GetUserPrincipalPreference(pid uint32) *friends_wiiu_types.PrincipalPreference {
	preference := friends_wiiu_types.NewPrincipalPreference()

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
