package database_wiiu

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/PretendoNetwork/friends/database"
	"github.com/PretendoNetwork/friends/globals"
	"github.com/PretendoNetwork/nex-go"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/friends-wiiu/types"
)

// GetUserFriendList returns a user's friend list
func GetUserFriendList(pid uint32) ([]*friends_wiiu_types.FriendInfo, error) {
	friendList := make([]*friends_wiiu_types.FriendInfo, 0)

	rows, err := database.Postgres.Query(`SELECT user2_pid, date FROM wiiu.friendships WHERE user1_pid=$1 AND active=true LIMIT 100`, pid)
	if err != nil {
		if err == sql.ErrNoRows {
			return friendList, database.ErrEmptyList
		} else {
			return friendList, err
		}
	}

	for rows.Next() {
		var friendPID uint32
		var date uint64
		rows.Scan(&friendPID, &date)

		friendInfo := friends_wiiu_types.NewFriendInfo()
		connectedUser := globals.ConnectedUsers[friendPID]
		var lastOnline *nex.DateTime

		if connectedUser != nil {
			// Online
			friendInfo.NNAInfo = connectedUser.NNAInfo
			friendInfo.Presence = connectedUser.PresenceV2

			if friendInfo.NNAInfo == nil || friendInfo.NNAInfo.PrincipalBasicInfo == nil {
				// TODO: Fix this
				globals.Logger.Error(fmt.Sprintf("User %d has friend %d with bad presence data", pid, friendPID))
				if friendInfo.NNAInfo == nil {
					globals.Logger.Error(fmt.Sprintf("%d friendInfo.NNAInfo is nil", friendPID))
				} else {
					globals.Logger.Error(fmt.Sprintf("%d friendInfo.NNAInfo.PrincipalBasicInfo is nil", friendPID))
				}

				continue
			}

			lastOnline = nex.NewDateTime(0)
			lastOnline.FromTimestamp(time.Now())
		} else {
			// Offline

			userData, err := globals.GetUserData(friendPID)
			if err != nil {
				globals.Logger.Critical(err.Error())
				continue
			}

			userInfo, err := GetUserInfoByPNIDData(userData)
			if err != nil {
				globals.Logger.Critical(err.Error())
				continue
			}

			friendInfo.NNAInfo = friends_wiiu_types.NewNNAInfo()

			friendInfo.NNAInfo.PrincipalBasicInfo = userInfo
			friendInfo.NNAInfo.Unknown1 = 0
			friendInfo.NNAInfo.Unknown2 = 0

			friendInfo.Presence = friends_wiiu_types.NewNintendoPresenceV2()
			friendInfo.Presence.ChangedFlags = 0
			friendInfo.Presence.Online = false
			friendInfo.Presence.GameKey = friends_wiiu_types.NewGameKey()
			friendInfo.Presence.GameKey.TitleID = 0
			friendInfo.Presence.GameKey.TitleVersion = 0
			friendInfo.Presence.Unknown1 = 0
			friendInfo.Presence.Message = ""
			friendInfo.Presence.Unknown2 = 0
			friendInfo.Presence.Unknown3 = 0
			friendInfo.Presence.GameServerID = 0
			friendInfo.Presence.Unknown4 = 0
			friendInfo.Presence.PID = 0
			friendInfo.Presence.GatheringID = 0
			friendInfo.Presence.ApplicationData = []byte{0x00}
			friendInfo.Presence.Unknown5 = 0
			friendInfo.Presence.Unknown6 = 0
			friendInfo.Presence.Unknown7 = 0

			var lastOnlineTime uint64
			err = database.Postgres.QueryRow(`SELECT last_online FROM wiiu.user_data WHERE pid=$1`, friendPID).Scan(&lastOnlineTime)
			if err != nil {
				globals.Logger.Critical(err.Error())
				continue
			}

			lastOnline = nex.NewDateTime(lastOnlineTime) // TODO: Change this
		}

		status, err := GetUserComment(friendPID)
		if err != nil {
			globals.Logger.Critical(err.Error())
			continue
		}

		friendInfo.Status = status

		friendInfo.BecameFriend = nex.NewDateTime(date)
		friendInfo.LastOnline = lastOnline
		friendInfo.Unknown = 0

		friendList = append(friendList, friendInfo)
	}

	return friendList, nil
}
