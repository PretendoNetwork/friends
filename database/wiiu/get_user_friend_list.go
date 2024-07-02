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

	rows, err := database.Manager.Query(`
	SELECT
		f.user2_pid, f.date,
		u.comment, u.comment_changed,
		u.last_online,
		bi.username, bi.unknown,
		ai.unknown1, ai.unknown2,
		mii.name, mii.unknown1, mii.unknown2, mii.data, mii.unknown_datetime
	FROM wiiu.friendships AS f
	INNER JOIN wiiu.user_data AS u ON u.pid = f.user2_pid
	INNER JOIN wiiu.principal_basic_info AS bi ON bi.pid = f.user2_pid
	INNER JOIN wiiu.network_account_info AS ai ON ai.pid = f.user2_pid
	INNER JOIN wiiu.mii AS mii ON mii.pid = f.user2_pid
	WHERE f.user1_pid=$1 AND f.active=true
	LIMIT 100
	`, pid)

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
		var lastOnlineTime uint64
		var commentContents string
		var commentChanged uint64 = 0
		var nnid string
		var unknown uint8
		var unknown1 uint8
		var unknown2 uint8
		var miiName string
		var miiUnknown1 uint8
		var miiUnknown2 uint8
		var miiData []byte
		var miiDatetime uint64

		err := rows.Scan(&friendPID, &date, &commentContents, &commentChanged, &lastOnlineTime, &nnid, &unknown, &unknown1, &unknown2, &miiName, &miiUnknown1, &miiUnknown2, &miiData, &miiDatetime)
		if err != nil {
			return nil, err
		}

		comment := friends_wiiu_types.NewComment()
		comment.Unknown = types.NewPrimitiveU8(0)
		comment.Contents = types.NewString(commentContents)
		comment.LastChanged = types.NewDateTime(commentChanged)

		mii := friends_wiiu_types.NewMiiV2()
		mii.Name = types.NewString(miiName)
		mii.Unknown1 = types.NewPrimitiveU8(miiUnknown1)
		mii.Unknown2 = types.NewPrimitiveU8(miiUnknown2)
		mii.MiiData = types.NewBuffer(miiData)
		mii.Datetime = types.NewDateTime(miiDatetime)

		principalBasicInfo := friends_wiiu_types.NewPrincipalBasicInfo()
		principalBasicInfo.PID = types.NewPID(uint64(friendPID))
		principalBasicInfo.NNID = types.NewString(nnid)
		principalBasicInfo.Unknown = types.NewPrimitiveU8(unknown)
		principalBasicInfo.Mii = mii

		nnaInfo := friends_wiiu_types.NewNNAInfo()
		nnaInfo.Unknown1 = types.NewPrimitiveU8(unknown1)
		nnaInfo.Unknown2 = types.NewPrimitiveU8(unknown2)
		nnaInfo.PrincipalBasicInfo = principalBasicInfo

		friendInfo := friends_wiiu_types.NewFriendInfo()
		friendInfo.NNAInfo = nnaInfo

		lastOnline := types.NewDateTime(0).Now()
		connectedUser, ok := globals.ConnectedUsers.Get(friendPID)
		if ok && connectedUser != nil {
			// * Online
			friendInfo.Presence = connectedUser.PresenceV2.Copy().(*friends_wiiu_types.NintendoPresenceV2)
		} else {
			// * Offline
			lastOnline = types.NewDateTime(lastOnlineTime) // TODO - Change this
		}

		friendInfo.Status = comment
		friendInfo.BecameFriend = types.NewDateTime(date)
		friendInfo.LastOnline = lastOnline
		friendInfo.Unknown = types.NewPrimitiveU64(0)

		friendList.Append(friendInfo)
	}

	return friendList, nil
}
