package database_wiiu

import (
	"github.com/PretendoNetwork/friends-secure/database"
)

func SetFriendRequestReceived(friendRequestID uint64) error {
	_, err := database.Postgres.Exec(`UPDATE wiiu.friend_requests SET received=true WHERE id=$1`, friendRequestID)

	if err != nil {
		return err
	}

	return nil
}
