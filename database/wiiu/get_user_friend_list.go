package database_wiiu

import (
	"fmt"
	"time"

	"github.com/PretendoNetwork/friends-secure/database"
	"github.com/PretendoNetwork/friends-secure/globals"
	"github.com/PretendoNetwork/nex-go"
	friends_wiiu "github.com/PretendoNetwork/nex-protocols-go/friends/wiiu"
	"github.com/gocql/gocql"
)

// Get a users friend list
func GetUserFriendList(pid uint32) []*friends_wiiu.FriendInfo {
	friendList := make([]*friends_wiiu.FriendInfo, 0)

	rows, err := database.Postgres.Query(`SELECT user2_pid, date FROM wiiu.friendships WHERE user1_pid=$1 AND active=true LIMIT 100`, pid)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return friendList
	}

	for rows.Next() {
		var friendPID uint32
		var date uint64
		rows.Scan(&friendPID, &date)

		friendInfo := friends_wiiu.NewFriendInfo()
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

			friendInfo.NNAInfo = friends_wiiu.NewNNAInfo()
			friendInfo.NNAInfo.PrincipalBasicInfo = GetUserInfoByPID(friendPID)
			friendInfo.NNAInfo.Unknown1 = 0
			friendInfo.NNAInfo.Unknown2 = 0

			friendInfo.Presence = friends_wiiu.NewNintendoPresenceV2()
			friendInfo.Presence.ChangedFlags = 0
			friendInfo.Presence.Online = false
			friendInfo.Presence.GameKey = friends_wiiu.NewGameKey()
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
			err := database.Postgres.QueryRow(`SELECT last_online FROM wiiu.user_data WHERE pid=$1`, friendPID).Scan(&lastOnlineTime)
			if err != nil {
				lastOnlineTime = nex.NewDateTime(0).Now()

				if err == gocql.ErrNotFound {
					globals.Logger.Error(err.Error())
				} else {
					globals.Logger.Critical(err.Error())
				}
			}

			lastOnline = nex.NewDateTime(lastOnlineTime) // TODO: Change this
		}

		friendInfo.Status = GetUserComment(friendPID)
		friendInfo.BecameFriend = nex.NewDateTime(date)
		friendInfo.LastOnline = lastOnline
		friendInfo.Unknown = 0

		friendList = append(friendList, friendInfo)
	}

	return friendList
}
