package database_wiiu

import (
	"database/sql"

	"github.com/PretendoNetwork/friends-secure/database"
	"github.com/PretendoNetwork/friends-secure/globals"
)

// Get a users outgoing friend request
func GetPIDsByFriendRequestID(friendRequestID uint64) (uint32, uint32) {
	var senderPID uint32
	var recipientPID uint32

	err := database.Postgres.QueryRow(`
	SELECT sender_pid, recipient_pid FROM wiiu.friend_requests WHERE id=$1
	`, friendRequestID).Scan(&senderPID, &recipientPID)
	if err != nil {
		if err == sql.ErrNoRows {
			globals.Logger.Warning(err.Error())
		} else {
			globals.Logger.Critical(err.Error())
		}
	}

	return senderPID, recipientPID
}
