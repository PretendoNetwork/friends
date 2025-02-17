package database_wiiu

import (
	"database/sql"

	"github.com/PretendoNetwork/friends/database"
)

// GetUserFriendPIDs returns a user's friend PIDs list
func GetUserFriendPIDs(pid uint32) ([]uint32, error) {
	pids := make([]uint32, 0)

	rows, err := database.Manager.Query(`SELECT user2_pid FROM wiiu.friendships WHERE user1_pid=$1 AND active=true LIMIT 100`, pid)
	if err != nil {
		if err == sql.ErrNoRows {
			return pids, database.ErrEmptyList
		} else {
			return pids, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		var pid uint32

		err := rows.Scan(&pid)
		if err != nil {
			return pids, err
		}

		pids = append(pids, pid)
	}

	return pids, nil
}
