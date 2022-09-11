package database

import (
	"database/sql"

	"github.com/PretendoNetwork/friends-secure/globals"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func GetUserPrincipalPreference(pid uint32) *nexproto.PrincipalPreference {
	preference := nexproto.NewPrincipalPreference()

	err := postgres.QueryRow(`SELECT show_online, show_current_game, block_friend_requests FROM wiiu.user_data WHERE pid=$1`, pid).Scan(&preference.ShowOnlinePresence, &preference.ShowCurrentTitle, &preference.BlockFriendRequests)
	if err != nil {
		if err == sql.ErrNoRows {
			globals.Logger.Warning(err.Error())
		} else {
			globals.Logger.Critical(err.Error())
		}
	}

	return preference
}
