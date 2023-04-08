package database_wiiu

import (
	"github.com/PretendoNetwork/friends-secure/database"
	"github.com/PretendoNetwork/friends-secure/globals"
)

// Remove a block from a user
func UnsetUserBlocked(user1_pid uint32, user2_pid uint32) {
	_, err := database.Postgres.Exec(`
		DELETE FROM wiiu.blocks WHERE blocker_pid=$1 AND blocked_pid=$2`, user1_pid, user2_pid)
	if err != nil {
		globals.Logger.Critical(err.Error())
	}
}
