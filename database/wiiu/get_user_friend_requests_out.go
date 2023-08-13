package database_wiiu

import (
	"database/sql"

	"github.com/PretendoNetwork/friends/database"
	"github.com/PretendoNetwork/friends/globals"
	"github.com/PretendoNetwork/nex-go"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/friends-wiiu/types"
)

// GetUserFriendRequestsOut returns the friend requests sent by a user
func GetUserFriendRequestsOut(pid uint32) ([]*friends_wiiu_types.FriendRequest, error) {
	friendRequestsOut := make([]*friends_wiiu_types.FriendRequest, 0)

	rows, err := database.Postgres.Query(`SELECT id, recipient_pid, sent_on, expires_on, message, received FROM wiiu.friend_requests WHERE sender_pid=$1 AND accepted=false`, pid)
	if err != nil {
		if err == sql.ErrNoRows {
			return friendRequestsOut, database.ErrEmptyList
		} else {
			return friendRequestsOut, err
		}
	}

	for rows.Next() {
		var id uint64
		var recipientPID uint32
		var sentOn uint64
		var expiresOn uint64
		var message string
		var received bool
		rows.Scan(&id, &recipientPID, &sentOn, &expiresOn, &message, &received)

		userData, err := globals.GetUserData(recipientPID)
		if err != nil {
			globals.Logger.Critical(err.Error())
			continue
		}

		userInfo, err := GetUserInfoByPNIDData(userData)
		if err != nil {
			globals.Logger.Critical(err.Error())
			continue
		}

		friendRequest := friends_wiiu_types.NewFriendRequest()

		friendRequest.PrincipalInfo = userInfo
		friendRequest.Message = friends_wiiu_types.NewFriendRequestMessage()
		friendRequest.Message.FriendRequestID = id
		friendRequest.Message.Received = received
		friendRequest.Message.Unknown2 = 1
		friendRequest.Message.Message = message
		friendRequest.Message.Unknown3 = 0
		friendRequest.Message.Unknown4 = ""
		friendRequest.Message.GameKey = friends_wiiu_types.NewGameKey()
		friendRequest.Message.GameKey.TitleID = 0
		friendRequest.Message.GameKey.TitleVersion = 0
		friendRequest.Message.Unknown5 = nex.NewDateTime(134222053376) // idk what this value means but its always this
		friendRequest.Message.ExpiresOn = nex.NewDateTime(expiresOn)
		friendRequest.SentOn = nex.NewDateTime(sentOn)

		friendRequestsOut = append(friendRequestsOut, friendRequest)
	}

	return friendRequestsOut, nil
}
