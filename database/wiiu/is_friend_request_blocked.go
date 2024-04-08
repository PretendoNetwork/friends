package database_wiiu

import (
	"github.com/PretendoNetwork/friends/database"
)

// IsFriendRequestBlocked determines if a requester PID has blocked a requested PID
func IsFriendRequestBlocked(requesterPID uint32, requestedPID uint32) (bool, error) {
	var found bool

	err := database.Postgres.QueryRow(`SELECT COUNT(*) FROM wiiu.blocks WHERE blocker_pid=$1 AND blocked_pid=$2 LIMIT 1`, requesterPID, requestedPID).Scan(&found)
	if err != nil {
		return false, err
	}

	return found, nil
}
