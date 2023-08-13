package database_wiiu

import (
	"github.com/PretendoNetwork/friends/database"
)

func DeleteFriendRequestAndReturnFriendPID(friendRequestID uint64) (uint32, error) {
	var recipientPID uint32

	err := database.Postgres.QueryRow(`SELECT recipient_pid FROM wiiu.friend_requests WHERE id=$1`, friendRequestID).Scan(&recipientPID)
	if err != nil {
		return 0, err
	}

	_, err = database.Postgres.Exec(`
		DELETE FROM wiiu.friend_requests WHERE id=$1`, friendRequestID)
	if err != nil {
		return 0, err
	}

	return recipientPID, nil
}
