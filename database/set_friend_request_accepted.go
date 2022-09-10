package database

func SetFriendRequestAccepted(friendRequestID uint64) {
	if err := cassandraClusterSession.Query(`UPDATE pretendo_friends.friend_requests SET accepted=true WHERE id=?`, friendRequestID).Exec(); err != nil {
		logger.Critical(err.Error())
	}
}
