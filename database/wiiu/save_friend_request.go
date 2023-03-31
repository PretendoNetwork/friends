package database_wiiu

import (
	"database/sql"

	"github.com/PretendoNetwork/friends-secure/database"
	"github.com/PretendoNetwork/friends-secure/globals"
)

func SaveFriendRequest(senderPID uint32, recipientPID uint32, sentTime uint64, expireTime uint64, message string) uint64 {
	var id uint64

	friendRequestBlocked := IsFriendRequestBlocked(recipientPID, senderPID)

	// Make sure we don't already have that friend request! If we do, give them the one we already have.
	err := database.Postgres.QueryRow(`SELECT id FROM wiiu.friend_requests WHERE sender_pid=$1 AND recipient_pid=$2`, senderPID, recipientPID).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		globals.Logger.Critical(err.Error())
		return 0
	} else if id != 0 {
		// If they aren't blocked, we want to unset the denied status on the previous request we have so that it appears again.
		if friendRequestBlocked {
			return id
		} else {
			UnsetFriendRequestDenied(id)
			return id
		}
	}

	err = database.Postgres.QueryRow(`
		INSERT INTO wiiu.friend_requests (sender_pid, recipient_pid, sent_on, expires_on, message, received, accepted, denied)
		VALUES ($1, $2, $3, $4, $5, false, false, $6) RETURNING id`, senderPID, recipientPID, sentTime, expireTime, message, friendRequestBlocked).Scan(&id)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return 0
	}

	return id
}
