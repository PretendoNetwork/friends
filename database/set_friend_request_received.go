package database

func SetFriendRequestReceived(friendRequestID uint64) {
	if err := cassandraClusterSession.Query(`UPDATE pretendo_friends.friend_requests SET received=true WHERE id=?`, friendRequestID).Exec(); err != nil {
		logger.Critical(err.Error())
	}
}
