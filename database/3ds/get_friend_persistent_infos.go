package database_3ds

import (
	"github.com/PretendoNetwork/friends/database"
	"github.com/PretendoNetwork/nex-go"
	friends_3ds_types "github.com/PretendoNetwork/nex-protocols-go/friends-3ds/types"
	"github.com/lib/pq"
)

// GetFriendPersistentInfos returns the persistent information of all friends
func GetFriendPersistentInfos(user1_pid uint32, pids []uint32) ([]*friends_3ds_types.FriendPersistentInfo, error) {
	persistentInfos := make([]*friends_3ds_types.FriendPersistentInfo, 0)

	rows, err := database.Postgres.Query(`
	SELECT pid, region, area, language, favorite_title, favorite_title_version, comment, comment_changed, last_online, mii_changed FROM "3ds".user_data WHERE pid=ANY($1::int[])`, pq.Array(pids))
	if err != nil {
		return persistentInfos, err
	}

	for rows.Next() {
		persistentInfo := friends_3ds_types.NewFriendPersistentInfo()

		gameKey := friends_3ds_types.NewGameKey()

		var pid uint32
		var lastOnlineTime uint64
		var msgUpdateTime uint64
		var miiModifiedAtTime uint64

		rows.Scan(
			&pid, &persistentInfo.Region, &persistentInfo.Area, &persistentInfo.Language,
			&gameKey.TitleID, &gameKey.TitleVersion, &persistentInfo.Message, &msgUpdateTime, &lastOnlineTime, &miiModifiedAtTime)

		persistentInfo.PID = nex.NewPID(pid)
		persistentInfo.MessageUpdatedAt = nex.NewDateTime(msgUpdateTime)
		persistentInfo.MiiModifiedAt = nex.NewDateTime(miiModifiedAtTime)
		persistentInfo.LastOnline = nex.NewDateTime(lastOnlineTime)
		persistentInfo.GameKey = gameKey
		persistentInfo.Platform = 2 // Always 3DS

		persistentInfos = append(persistentInfos, persistentInfo)
	}

	return persistentInfos, nil
}
