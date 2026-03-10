package database_wiiu

import (
	"database/sql"

	"github.com/PretendoNetwork/friends/database"
)

// SaveFriendRequest registers a new friend request
func SaveFriendRequest(senderPID uint32, recipientPID uint32, sentTime uint64, expireTime uint64, message string) (uint64, error) {
	var id uint64

	friendRequestBlocked, err := IsFriendRequestBlocked(recipientPID, senderPID)
	if err != nil {
		return 0, err
	}

	// Make sure we don't already have that friend request! If we do, give them the one we already have.
	row, err := database.Manager.QueryRow(`SELECT id FROM wiiu.friend_requests WHERE sender_pid=$1 AND recipient_pid=$2`, senderPID, recipientPID)
	if err != nil {
		return 0, err
	}

	err = row.Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	} else if id != 0 {
		// A previous friend request already exists. Reset its state so it
		// appears as a fresh incoming request (handles the case where a
		// friendship was accepted and later removed).
		if friendRequestBlocked {
			return id, nil
		}

		_, err = database.Manager.Exec(
			`UPDATE wiiu.friend_requests
			 SET denied=false, accepted=false, received=false,
			     sent_on=$1, expires_on=$2, message=$3
			 WHERE id=$4`,
			sentTime, expireTime, message, id)
		if err != nil {
			return 0, err
		}

		return id, nil
	}

	row, err = database.Manager.QueryRow(`
		INSERT INTO wiiu.friend_requests (sender_pid, recipient_pid, sent_on, expires_on, message, received, accepted, denied)
		VALUES ($1, $2, $3, $4, $5, false, false, $6) RETURNING id`, senderPID, recipientPID, sentTime, expireTime, message, friendRequestBlocked)
	if err != nil {
		return 0, err
	}

	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
