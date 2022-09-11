package database

import "github.com/PretendoNetwork/friends-secure/globals"

func SetFriendRequestReceived(friendRequestID uint64) {
	_, err := postgres.Exec(`UPDATE wiiu.friend_requests SET received=true WHERE id=$1`, friendRequestID)

	if err != nil {
		globals.Logger.Critical(err.Error())
	}
}
