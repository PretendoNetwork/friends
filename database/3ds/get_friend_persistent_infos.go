package database_3ds

import (
	"github.com/PretendoNetwork/friends/database"
	"github.com/PretendoNetwork/nex-go/v2/types"
	friends_3ds_types "github.com/PretendoNetwork/nex-protocols-go/v2/friends-3ds/types"
	"github.com/lib/pq"
)

// GetFriendPersistentInfos returns the persistent information of all friends
func GetFriendPersistentInfos(user1_pid uint32, pids []uint32) (*types.List[*friends_3ds_types.FriendPersistentInfo], error) {
	persistentInfos := types.NewList[*friends_3ds_types.FriendPersistentInfo]()
	persistentInfos.Type = friends_3ds_types.NewFriendPersistentInfo()

	rows, err := database.Postgres.Query(`
	SELECT pid, region, area, language, favorite_title, favorite_title_version, comment, comment_changed, last_online, mii_changed FROM "3ds".user_data WHERE pid=ANY($1::int[])`, pq.Array(pids))
	if err != nil {
		return persistentInfos, err
	}

	for rows.Next() {
		persistentInfo := friends_3ds_types.NewFriendPersistentInfo()

		gameKey := friends_3ds_types.NewGameKey()

		var pid uint32
		var region uint8
		var area uint8
		var language uint8
		var titleID uint64
		var titleVersion uint16
		var message string
		var lastOnlineTime uint64
		var msgUpdateTime uint64
		var miiModifiedAtTime uint64

		// * This is allowed to error for now.
		// * Some of these fields are optional, and the DB doesn't have defaults
		rows.Scan(
			&pid,
			&region,
			&area,
			&language,
			&titleID,
			&titleVersion,
			&message,
			&msgUpdateTime,
			&lastOnlineTime,
			&miiModifiedAtTime,
		)

		gameKey.TitleID = types.NewPrimitiveU64(titleID)
		gameKey.TitleVersion = types.NewPrimitiveU16(titleVersion)

		persistentInfo.PID = types.NewPID(uint64(pid))
		persistentInfo.Region = types.NewPrimitiveU8(region)
		persistentInfo.Country = types.NewPrimitiveU8(0) // TODO - What is this?
		persistentInfo.Area = types.NewPrimitiveU8(area)
		persistentInfo.Language = types.NewPrimitiveU8(language)
		persistentInfo.Platform = types.NewPrimitiveU8(2) // * Always 3DS
		persistentInfo.GameKey = gameKey
		persistentInfo.Message = types.NewString(message)
		persistentInfo.MessageUpdatedAt = types.NewDateTime(msgUpdateTime)
		persistentInfo.MiiModifiedAt = types.NewDateTime(miiModifiedAtTime)
		persistentInfo.LastOnline = types.NewDateTime(lastOnlineTime)

		persistentInfos.Append(persistentInfo)
	}

	return persistentInfos, nil
}
