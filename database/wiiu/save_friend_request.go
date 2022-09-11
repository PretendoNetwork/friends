package database_wiiu

import (
	"github.com/PretendoNetwork/friends-secure/database"
	"github.com/PretendoNetwork/friends-secure/globals"
)

func SaveFriendRequest(senderPID uint32, recipientPID uint32, sentTime uint64, expireTime uint64, message string) uint64 {
	var id uint64
	err := database.Postgres.QueryRow(`
		INSERT INTO wiiu.friend_requests (sender_pid, recipient_pid, sent_on, expires_on, message, received, accepted, denied)
		VALUES ($1, $2, $3, $4, $5, false, false, false) RETURNING id`, senderPID, recipientPID, sentTime, expireTime, message).Scan(&id)

	if err != nil {
		globals.Logger.Critical(err.Error())
		return 0
	}

	return id
}
