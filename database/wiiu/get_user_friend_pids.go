package database_wiiu

import (
	"github.com/PretendoNetwork/friends-secure/database"
	"github.com/PretendoNetwork/friends-secure/globals"
)

// Get a users friend PIDs list
func GetUserFriendPIDs(pid uint32) []uint32 {
	pids := make([]uint32, 0)

	rows, err := database.Postgres.Query(`SELECT user2_pid FROM wiiu.friendships WHERE user1_pid=$1 AND active=true LIMIT 100`, pid)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return pids
	}
	defer rows.Close()

	for rows.Next() {
		var pid uint32
		rows.Scan(&pid)

		pids = append(pids, pid)
	}

	return pids
}
