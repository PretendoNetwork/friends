package database_wiiu

import (
	"github.com/PretendoNetwork/friends-secure/database"
	"github.com/PretendoNetwork/friends-secure/globals"
)

func IsFriendRequestBlocked(requesterPID uint32, requestedPID uint32) bool {
	var found bool
	err := database.Postgres.QueryRow(`SELECT COUNT(*) FROM wiiu.blocks WHERE blocker_pid=$1 AND blocked_pid=$2 LIMIT 1`, requesterPID, requestedPID).Scan(&found)
	if err != nil {
		globals.Logger.Critical(err.Error())
	}

	return found
}
