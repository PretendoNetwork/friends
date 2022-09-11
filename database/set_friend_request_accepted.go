package database

import "github.com/PretendoNetwork/friends-secure/globals"

func SetFriendRequestAccepted(friendRequestID uint64) {
	_, err := postgres.Exec(`UPDATE wiiu.friend_requests SET accepted=true WHERE id=$1`, friendRequestID)

	if err != nil {
		globals.Logger.Critical(err.Error())
	}
}
