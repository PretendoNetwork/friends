package database_wiiu

import (
	"github.com/PretendoNetwork/friends-secure/database"
	"github.com/PretendoNetwork/friends-secure/globals"
)

func UnsetFriendRequestDenied(friendRequestID uint64) {
	_, err := database.Postgres.Exec(`UPDATE wiiu.friend_requests SET denied=false WHERE id=$1`, friendRequestID)

	if err != nil {
		globals.Logger.Critical(err.Error())
	}
}
