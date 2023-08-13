package database_wiiu

import (
	"github.com/PretendoNetwork/friends/database"
)

func SetFriendRequestAccepted(friendRequestID uint64) error {
	_, err := database.Postgres.Exec(`UPDATE wiiu.friend_requests SET accepted=true WHERE id=$1`, friendRequestID)
	if err != nil {
		return err
	}

	return nil
}
