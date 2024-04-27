package database_wiiu

import (
	"database/sql"

	"github.com/PretendoNetwork/friends/database"
)

// GetPIDsByFriendRequestID returns the users outgoing friend request
func GetPIDsByFriendRequestID(friendRequestID uint64) (uint32, uint32, error) {
	var senderPID uint32
	var recipientPID uint32

	row, err := database.Manager.QueryRow(`SELECT sender_pid, recipient_pid FROM wiiu.friend_requests WHERE id=$1`, friendRequestID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, 0, database.ErrFriendRequestNotFound
		} else {
			return 0, 0, err
		}
	}

	err = row.Scan(&senderPID, &recipientPID)
	if err != nil {
		return 0, 0, err
	}

	return senderPID, recipientPID, nil
}
