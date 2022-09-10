package database

import "github.com/gocql/gocql"

func IsFriendRequestBlocked(requesterPID uint32, requestedPID uint32) bool {
	if err := cassandraClusterSession.Query(`SELECT id FROM pretendo_friends.blocks WHERE blocker_pid=? AND blocked_pid=? LIMIT 1 ALLOW FILTERING`, requestedPID, requesterPID).Scan(); err != nil {
		if err == gocql.ErrNotFound {
			// Assume no block record was found
			return false
		}

		// TODO: Error handling
	}

	// Assume a block record was found
	return true
}
