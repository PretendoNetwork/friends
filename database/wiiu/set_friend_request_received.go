package database_wiiu

import (
	"github.com/PretendoNetwork/friends/database"
)

// SetFriendRequestReceived marks a friend request as received
func SetFriendRequestReceived(friendRequestID uint64) error {
	result, err := database.Postgres.Exec(`UPDATE wiiu.friend_requests SET received=true WHERE id=$1`, friendRequestID)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return database.ErrFriendRequestNotFound
	}

	return nil
}
