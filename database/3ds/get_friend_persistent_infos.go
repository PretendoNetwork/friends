package database_3ds

import (
	"database/sql"
	"time"

	"github.com/PretendoNetwork/friends-secure/database"
	"github.com/PretendoNetwork/friends-secure/globals"
	"github.com/PretendoNetwork/nex-go"
	friends_3ds "github.com/PretendoNetwork/nex-protocols-go/friends/3ds"
)

// Get a friend's persistent information
func GetFriendPersistentInfos(user1_pid uint32, pids []uint32) []*friends_3ds.FriendPersistentInfo {
	persistentInfos := make([]*friends_3ds.FriendPersistentInfo, 0)

	rows, err := database.Postgres.Query(`
	SELECT pid, region, area, language, favorite_title, favorite_title_version, comment, comment_changed, last_online FROM "3ds".user_data WHERE pid IN ($1)`, database.PIDArrayToString(pids))
	if err != nil {
		if err == sql.ErrNoRows {
			globals.Logger.Warning(err.Error())
		} else {
			globals.Logger.Critical(err.Error())
		}
	}

	for rows.Next() {
		persistentInfo := friends_3ds.NewFriendPersistentInfo()

		gameKey := friends_3ds.NewGameKey()

		var lastOnlineTime uint64
		var msgUpdateTime uint64
		var friendedAtTime uint64

		rows.Scan(
			&persistentInfo.PID, &persistentInfo.Region, &persistentInfo.Area, &persistentInfo.Language,
			&gameKey.TitleID, &gameKey.TitleVersion, &persistentInfo.Message, &msgUpdateTime, &lastOnlineTime)

		err = database.Postgres.QueryRow(`
			SELECT date FROM "3ds".friendships WHERE user1_pid=$1 AND user2_pid=$2 AND type=0 LIMIT 1`, user1_pid, persistentInfo.PID).Scan(&friendedAtTime)
		if err != nil {
			if err == sql.ErrNoRows {
				friendedAtTime = uint64(time.Now().Unix())
			} else {
				globals.Logger.Critical(err.Error())
			}
		}

		persistentInfo.MessageUpdatedAt = nex.NewDateTime(msgUpdateTime)
		persistentInfo.FriendedAt = nex.NewDateTime(friendedAtTime)
		persistentInfo.LastOnline = nex.NewDateTime(lastOnlineTime)
		persistentInfo.GameKey = gameKey
		persistentInfo.Platform = 2 // Always 3DS

		persistentInfos = append(persistentInfos, persistentInfo)
	}

	return persistentInfos
}
