package database

func SaveFriendRequest(friendRequestID uint64, senderPID uint32, recipientPID uint32, sentTime uint64, expireTime uint64, message string) {
	if err := cassandraClusterSession.Query(`INSERT INTO pretendo_friends.friend_requests (id, sender_pid, recipient_pid, sent_on, expires_on, message, received, accepted, denied) VALUES (?, ?, ?, ?, ?, ?, false, false, false) IF NOT EXISTS`, friendRequestID, senderPID, recipientPID, sentTime, expireTime, message).Exec(); err != nil {
		logger.Critical(err.Error())
	}
}
