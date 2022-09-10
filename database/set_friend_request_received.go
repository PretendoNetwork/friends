package database

import "github.com/PretendoNetwork/friends-secure/globals"

func SetFriendRequestReceived(friendRequestID uint64) {
	if err := cassandraClusterSession.Query(`UPDATE pretendo_friends.friend_requests SET received=true WHERE id=?`, friendRequestID).Exec(); err != nil {
		globals.Logger.Critical(err.Error())
	}
}
