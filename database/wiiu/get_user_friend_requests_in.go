package database_wiiu

import (
	"database/sql"
	"time"

	"github.com/PretendoNetwork/friends/database"
	"github.com/PretendoNetwork/friends/utility"
	"github.com/PretendoNetwork/nex-go/v2/types"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/v2/friends-wiiu/types"
)

// GetUserFriendRequestsIn returns the friend requests received by a user
func GetUserFriendRequestsIn(pid uint32) (*types.List[*friends_wiiu_types.FriendRequest], error) {
	friendRequests := types.NewList[*friends_wiiu_types.FriendRequest]()
	friendRequests.Type = friends_wiiu_types.NewFriendRequest()

	rows, err := database.Postgres.Query(`SELECT id, sender_pid, sent_on, expires_on, message, received FROM wiiu.friend_requests WHERE recipient_pid=$1 AND accepted=false AND denied=false`, pid)
	if err != nil {
		if err == sql.ErrNoRows {
			return friendRequests, database.ErrEmptyList
		} else {
			return friendRequests, err
		}
	}

	for rows.Next() {
		var id uint64
		var senderPID uint32
		var sentOn uint64
		var expiresOn uint64
		var message string
		var received bool

		err := rows.Scan(&id, &senderPID, &sentOn, &expiresOn, &message, &received)
		if err != nil {
			return friendRequests, err
		}

		userInfo, err := utility.GetUserInfoByPID(senderPID)
		if err != nil {
			return nil, err
		}

		friendRequest := friends_wiiu_types.NewFriendRequest()

		friendRequest.PrincipalInfo = userInfo
		friendRequest.Message = friends_wiiu_types.NewFriendRequestMessage()
		friendRequest.Message.FriendRequestID = types.NewPrimitiveU64(id)
		friendRequest.Message.Received = types.NewPrimitiveBool(received)
		friendRequest.Message.Unknown2 = types.NewPrimitiveU8(1)
		friendRequest.Message.Message = types.NewString(message)
		friendRequest.Message.Unknown3 = types.NewPrimitiveU8(0)
		friendRequest.Message.Unknown4 = types.NewString("")
		friendRequest.Message.GameKey = friends_wiiu_types.NewGameKey()
		friendRequest.Message.GameKey.TitleID = types.NewPrimitiveU64(0)
		friendRequest.Message.GameKey.TitleVersion = types.NewPrimitiveU16(0)
		friendRequest.Message.Unknown5 = types.NewDateTime(134222053376) // * idk what this value means but its always this
		friendRequest.Message.ExpiresOn = types.NewDateTime(expiresOn)
		friendRequest.SentOn = types.NewDateTime(sentOn)

		// * Filter out expired requests
		if friendRequest.Message.ExpiresOn.Standard().After(time.Now()) {
			friendRequests.Append(friendRequest)
		}
	}

	return friendRequests, nil
}
