package database_wiiu

import (
	"github.com/PretendoNetwork/friends-secure/database"
	"github.com/PretendoNetwork/friends-secure/globals"
	"github.com/PretendoNetwork/nex-go"
	friends_wiiu "github.com/PretendoNetwork/nex-protocols-go/friends/wiiu"
)

// Get a users sent friend requests
func GetUserFriendRequestsOut(pid uint32) []*friends_wiiu.FriendRequest {
	friendRequestsOut := make([]*friends_wiiu.FriendRequest, 0)

	rows, err := database.Postgres.Query(`SELECT id, recipient_pid, sent_on, expires_on, message, received FROM wiiu.friend_requests WHERE sender_pid=$1 AND accepted=false`, pid)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return friendRequestsOut
	}

	for rows.Next() {
		var id uint64
		var recipientPID uint32
		var sentOn uint64
		var expiresOn uint64
		var message string
		var received bool
		rows.Scan(&id, &recipientPID, &sentOn, &expiresOn, &message, &received)

		friendRequest := friends_wiiu.NewFriendRequest()

		friendRequest.PrincipalInfo = GetUserInfoByPID(recipientPID)

		friendRequest.Message = friends_wiiu.NewFriendRequestMessage()
		friendRequest.Message.FriendRequestID = id
		friendRequest.Message.Received = received
		friendRequest.Message.Unknown2 = 1
		friendRequest.Message.Message = message
		friendRequest.Message.Unknown3 = 0
		friendRequest.Message.Unknown4 = ""
		friendRequest.Message.GameKey = friends_wiiu.NewGameKey()
		friendRequest.Message.GameKey.TitleID = 0
		friendRequest.Message.GameKey.TitleVersion = 0
		friendRequest.Message.Unknown5 = nex.NewDateTime(134222053376) // idk what this value means but its always this
		friendRequest.Message.ExpiresOn = nex.NewDateTime(expiresOn)
		friendRequest.SentOn = nex.NewDateTime(sentOn)

		friendRequestsOut = append(friendRequestsOut, friendRequest)
	}

	return friendRequestsOut
}
