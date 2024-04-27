package database_wiiu

import (
	"database/sql"
	"fmt"

	"github.com/PretendoNetwork/friends/database"
	"github.com/PretendoNetwork/friends/globals"
	"github.com/PretendoNetwork/friends/utility"
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
		var lastOnline *types.DateTime

		if ok && connectedUser != nil {
			// * Online
			friendInfo.NNAInfo = connectedUser.NNAInfo
			friendInfo.Presence = connectedUser.PresenceV2

			if friendInfo.NNAInfo == nil || friendInfo.NNAInfo.PrincipalBasicInfo == nil {
				// TODO - Fix this
				globals.Logger.Error(fmt.Sprintf("User %d has friend %d with bad presence data", pid, friendPID))
				if friendInfo.NNAInfo == nil {
					globals.Logger.Error(fmt.Sprintf("%d friendInfo.NNAInfo is nil", friendPID))
				} else {
					globals.Logger.Error(fmt.Sprintf("%d friendInfo.NNAInfo.PrincipalBasicInfo is nil", friendPID))
				}

				continue
			}

			lastOnline = types.NewDateTime(0).Now()
		} else {
			// * Offline

			userInfo, err := utility.GetUserInfoByPID(friendPID)
			if err != nil {
				return nil, err
			}

			friendInfo.NNAInfo = friends_wiiu_types.NewNNAInfo()

			friendInfo.NNAInfo.PrincipalBasicInfo = userInfo
			friendInfo.NNAInfo.Unknown1 = types.NewPrimitiveU8(0)
			friendInfo.NNAInfo.Unknown2 = types.NewPrimitiveU8(0)

			friendInfo.Presence = friends_wiiu_types.NewNintendoPresenceV2()
			friendInfo.Presence.ChangedFlags = types.NewPrimitiveU32(0)
			friendInfo.Presence.Online = types.NewPrimitiveBool(false)
			friendInfo.Presence.GameKey = friends_wiiu_types.NewGameKey()
			friendInfo.Presence.GameKey.TitleID = types.NewPrimitiveU64(0)
			friendInfo.Presence.GameKey.TitleVersion = types.NewPrimitiveU16(0)
			friendInfo.Presence.Unknown1 = types.NewPrimitiveU8(0)
			friendInfo.Presence.Message = types.NewString("")
			friendInfo.Presence.Unknown2 = types.NewPrimitiveU32(0)
			friendInfo.Presence.Unknown3 = types.NewPrimitiveU8(0)
			friendInfo.Presence.GameServerID = types.NewPrimitiveU32(0)
			friendInfo.Presence.Unknown4 = types.NewPrimitiveU32(0)
			friendInfo.Presence.PID = types.NewPID(0)
			friendInfo.Presence.GatheringID = types.NewPrimitiveU32(0)
			friendInfo.Presence.ApplicationData = types.NewBuffer([]byte{})
			friendInfo.Presence.Unknown5 = types.NewPrimitiveU8(0)
			friendInfo.Presence.Unknown6 = types.NewPrimitiveU8(0)
			friendInfo.Presence.Unknown7 = types.NewPrimitiveU8(0)

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
