package database_wiiu

import (
	"github.com/PretendoNetwork/friends/database"
)

// SetFriendRequestDenied marks a friend request as denied
func SetFriendRequestDenied(friendRequestID uint64) error {
	result, err := database.Manager.Exec(`UPDATE wiiu.friend_requests SET denied=true WHERE id=$1`, friendRequestID)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return database.ErrFriendRequestNotFound
	}

	return nil
}
