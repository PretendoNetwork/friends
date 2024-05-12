package database_wiiu

import (
	"database/sql"

	"github.com/PretendoNetwork/friends/database"
	"github.com/PretendoNetwork/friends/globals"
	"github.com/PretendoNetwork/nex-go/v2/types"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/v2/friends-wiiu/types"
)

// GetUserFriendList returns a user's friend list
func GetUserFriendList(pid uint32) (*types.List[*friends_wiiu_types.FriendInfo], error) {
	friendList := types.NewList[*friends_wiiu_types.FriendInfo]()
	friendList.Type = friends_wiiu_types.NewFriendInfo()

	rows, err := database.Manager.Query(`SELECT user2_pid, date FROM wiiu.friendships WHERE user1_pid=$1 AND active=true LIMIT 100`, pid)
	if err != nil {
		if err == sql.ErrNoRows {
			return friendList, database.ErrEmptyList
		} else {
			return friendList, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		var friendPID uint32
		var date uint64

		err := rows.Scan(&friendPID, &date)
		if err != nil {
			return nil, err
		}

		friendInfo := friends_wiiu_types.NewFriendInfo()
		connectedUser, ok := globals.ConnectedUsers.Get(friendPID)
		lastOnline := types.NewDateTime(0).Now()

		friendInfo.NNAInfo, err = GetUserNetworkAccountInfo(friendPID)
		if err != nil {
			return nil, err
		}

		if ok && connectedUser != nil {
			// * Online
			friendInfo.Presence = connectedUser.PresenceV2
		} else {
			// * Offline
			var lastOnlineTime uint64
			row, err := database.Manager.QueryRow(`SELECT last_online FROM wiiu.user_data WHERE pid=$1`, friendPID)
			if err != nil {
				return nil, err
			}

			err = row.Scan(&lastOnlineTime)
			if err != nil {
				return nil, err
			}

			lastOnline = types.NewDateTime(lastOnlineTime) // TODO - Change this
		}

		status, err := GetUserComment(friendPID)
		if err != nil {
			return nil, err
		}

		friendInfo.Status = status

		friendInfo.BecameFriend = types.NewDateTime(date)
		friendInfo.LastOnline = lastOnline
		friendInfo.Unknown = types.NewPrimitiveU64(0)

		friendList.Append(friendInfo)
	}

	return friendList, nil
}
