package database_wiiu

import (
	"github.com/PretendoNetwork/friends/database"
)

func UnsetFriendRequestDenied(friendRequestID uint64) error {
	_, err := database.Postgres.Exec(`UPDATE wiiu.friend_requests SET denied=false WHERE id=$1`, friendRequestID)

	if err != nil {
		return err
	}

	return nil
}
