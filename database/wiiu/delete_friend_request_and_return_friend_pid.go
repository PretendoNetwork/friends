package database_wiiu

import (
	"github.com/PretendoNetwork/friends-secure/database"
	"github.com/PretendoNetwork/friends-secure/globals"
)

func DeleteFriendRequestAndReturnFriendPID(friendRequestID uint64) uint32 {
	var recipientPID uint32

	err := database.Postgres.QueryRow(`SELECT recipient_pid FROM wiiu.friend_requests WHERE id=$1`, friendRequestID).Scan(&recipientPID)
	if err != nil {
		globals.Logger.Critical(err.Error())
	}

	_, err = database.Postgres.Exec(`
		DELETE FROM wiiu.friend_requests WHERE id=$1`, friendRequestID)
	if err != nil {
		globals.Logger.Critical(err.Error())
	}

	return recipientPID
}
