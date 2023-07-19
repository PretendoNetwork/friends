package database_3ds

import (
	"database/sql"
	"time"

	"github.com/PretendoNetwork/friends-secure/database"
	"github.com/PretendoNetwork/friends-secure/globals"
	"github.com/PretendoNetwork/nex-go"
	friends_3ds_types "github.com/PretendoNetwork/nex-protocols-go/friends-3ds/types"
)

// Get a friend's mii
func GetFriendMiis(pids []uint32) []*friends_3ds_types.FriendMii {
	friendMiis := make([]*friends_3ds_types.FriendMii, 0)

	rows, err := database.Postgres.Query(`
	SELECT pid, mii_name, mii_data FROM "3ds".user_data WHERE pid IN ($1)`, database.PIDArrayToString(pids))
	if err != nil {
		if err == sql.ErrNoRows {
			globals.Logger.Warning(err.Error())
		} else {
			globals.Logger.Critical(err.Error())
		}
	}

	changedTime := nex.NewDateTime(0)
	changedTime.FromTimestamp(time.Now())

	for rows.Next() {
		var pid uint32

		mii := friends_3ds_types.NewMii()
		mii.Unknown2 = false
		mii.Unknown3 = 0

		rows.Scan(&pid, &mii.Name, &mii.MiiData)

		friendMii := friends_3ds_types.NewFriendMii()
		friendMii.PID = pid
		friendMii.Mii = mii
		friendMii.ModifiedAt = changedTime

		friendMiis = append(friendMiis, friendMii)
	}

	return friendMiis
}
