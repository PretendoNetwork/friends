package database_wiiu

import (
	"database/sql"

	"github.com/PretendoNetwork/friends/database"
)

// DeleteFriendRequestAndReturnFriendPID deletes a given friend request and returns the friend's PID
func DeleteFriendRequestAndReturnFriendPID(friendRequestID uint64) (uint32, error) {
	var recipientPID uint32

	row, err := database.Manager.QueryRow(`SELECT recipient_pid FROM wiiu.friend_requests WHERE id=$1`, friendRequestID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, database.ErrFriendRequestNotFound
		} else {
			return 0, err
		}
	}

	err = row.Scan(&recipientPID)
	if err != nil {
		return 0, err
	}

	result, err := database.Manager.Exec(`
		DELETE FROM wiiu.friend_requests WHERE id=$1`, friendRequestID)
	if err != nil {
		return 0, err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return 0, database.ErrFriendRequestNotFound
	}

	return recipientPID, nil
}
