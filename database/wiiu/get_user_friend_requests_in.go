package database_wiiu

import (
	"database/sql"
	"time"

	"github.com/PretendoNetwork/friends/database"
	"github.com/PretendoNetwork/friends/utility"
	"github.com/PretendoNetwork/nex-go"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/friends-wiiu/types"
)

// GetUserFriendRequestsIn returns the friend requests received by a user
func GetUserFriendRequestsIn(pid uint32) ([]*friends_wiiu_types.FriendRequest, error) {
	friendRequestsIn := make([]*friends_wiiu_types.FriendRequest, 0)

	rows, err := database.Postgres.Query(`SELECT id, sender_pid, sent_on, expires_on, message, received FROM wiiu.friend_requests WHERE recipient_pid=$1 AND accepted=false AND denied=false`, pid)
	if err != nil {
		if err == sql.ErrNoRows {
			return friendRequestsIn, database.ErrEmptyList
		} else {
			return friendRequestsIn, err
		}
	}

	for rows.Next() {
		var id uint64
		var senderPID uint32
		var sentOn uint64
		var expiresOn uint64
		var message string
		var received bool
		rows.Scan(&id, &senderPID, &sentOn, &expiresOn, &message, &received)

		userInfo, err := utility.GetUserInfoByPID(senderPID)
		if err != nil {
			return make([]*friends_wiiu_types.FriendRequest, 0), err
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

		// * Filter out expired requests
		if friendRequest.Message.ExpiresOn.Standard().After(time.Now()) {
			friendRequestsIn = append(friendRequestsIn, friendRequest)
		}
	}

	return friendRequestsIn, nil
}
