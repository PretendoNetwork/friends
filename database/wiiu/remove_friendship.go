package database_wiiu

import (
	"database/sql"

	"github.com/PretendoNetwork/friends/database"
)

// RemoveFriendship removes a user's friend relationship
func RemoveFriendship(user1_pid uint32, user2_pid uint32) error {
	var completed bool

	row, err := database.Manager.QueryRow(`
		DELETE FROM wiiu.friendships WHERE user1_pid=$1 AND user2_pid=$2 RETURNING completed`, user1_pid, user2_pid)
	if err != nil {
		return err
	}

	err = row.Scan(&completed)
	if err != nil {
		if err == sql.ErrNoRows {
			return database.ErrFriendshipNotFound
		} else {
			return err
		}
	}

	if !completed {
		// * Friendship is provisional. Remove the associated friend request with it.
		friendRequestID, _, err := CheckExistingFriendRequestByPIDs(user1_pid, user2_pid)
		if err != nil {
			if err == sql.ErrNoRows {
				return database.ErrFriendRequestNotFound
			} else {
				return err
			}
		}

		err = SetFriendRequestDenied(friendRequestID)
		if err != nil {
			return err
		}

		return nil
	}

	_, err = database.Manager.Exec(`
		UPDATE wiiu.friendships SET active=false WHERE user1_pid=$1 AND user2_pid=$2`, user2_pid, user1_pid)
	if err != nil {
		return err
	}

	return nil
}
