package database_wiiu

import (
	"database/sql"

	"github.com/PretendoNetwork/friends/database"
)

// SaveProvisionalFriend creates a friend request and an incomplete friendship
func SaveProvisionalFriend(senderPID uint32, recipientPID uint32, sentTime uint64) (uint64, error) {
	var id uint64

	friendRequestBlocked, err := IsFriendRequestBlocked(recipientPID, senderPID)
	if err != nil {
		return 0, err
	}

	// Check for an existing friend request between the two users
	row, err := database.Manager.QueryRow(`SELECT id FROM wiiu.friend_requests WHERE sender_pid=$1 AND recipient_pid=$2`, senderPID, recipientPID)
	if err != nil {
		return 0, err
	}

	err = row.Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	} else if id != 0 {
		if friendRequestBlocked {
			return id, nil
		}

		// Reset status flags and update timestamps for resent requests
		_, err = database.Manager.Exec(`
			UPDATE wiiu.friend_requests 
			SET denied = false, 
			    accepted = false, 
			    sent_on = $1,
			    expires_on = 134222053376,
				provisional = true
			WHERE id = $2`,
			sentTime, id)

		if err != nil {
			return 0, err
		}
	}

	// Create a new provisional friend request if none exists
	row, err = database.Manager.QueryRow(`
		INSERT INTO wiiu.friend_requests (sender_pid, recipient_pid, sent_on, expires_on, message, received, accepted, denied, provisional)
		VALUES ($1, $2, $3, 134222053376, '', false, false, $4, true) RETURNING id`, senderPID, recipientPID, sentTime, friendRequestBlocked)
	if err != nil {
		return 0, err
	}

	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}

	// Create a provisional/incomplete friendship

	_, err = database.Manager.Exec(`
		INSERT INTO wiiu.friendships (user1_pid, user2_pid, date, active, completed)
		VALUES ($1, $2, $3, true, false)
		ON CONFLICT (user1_pid, user2_pid)
		DO UPDATE SET
		date = $3,
		active = true,
		completed = false`, senderPID, recipientPID, uint64(sentTime))
	if err != nil {
		return 0, err
	}

	return id, err
}
