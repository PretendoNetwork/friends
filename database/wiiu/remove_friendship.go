package database_wiiu

import (
	"github.com/PretendoNetwork/friends/database"
)

// RemoveFriendship removes a user's friend relationship
func RemoveFriendship(user1_pid uint32, user2_pid uint32) error {
	result, err := database.Manager.Exec(`
		DELETE FROM wiiu.friendships WHERE user1_pid=$1 AND user2_pid=$2`, user1_pid, user2_pid)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return database.ErrFriendshipNotFound
	}

	_, err = database.Manager.Exec(`
		UPDATE wiiu.friendships SET active=false WHERE user1_pid=$1 AND user2_pid=$2`, user2_pid, user1_pid)
	if err != nil {
		return err
	}

	return nil
}
