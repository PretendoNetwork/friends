package database_3ds

import (
	"github.com/PretendoNetwork/friends/database"
	"github.com/PretendoNetwork/nex-go"
	friends_3ds_types "github.com/PretendoNetwork/nex-protocols-go/friends-3ds/types"
	"github.com/lib/pq"
)

// GetFriendMiis returns the Mii of all friends
func GetFriendMiis(pids []uint32) ([]*friends_3ds_types.FriendMii, error) {
	friendMiis := make([]*friends_3ds_types.FriendMii, 0)

	rows, err := database.Postgres.Query(`
	SELECT pid, mii_name, mii_data FROM "3ds".user_data WHERE pid=ANY($1::int[])`, pq.Array(pids))
	if err != nil {
		return friendMiis, err
	}

	changedTime := nex.NewDateTime(0).Now()

	for rows.Next() {
		var pid uint32

		mii := friends_3ds_types.NewMii()
		mii.Unknown2 = false
		mii.Unknown3 = 0

		rows.Scan(&pid, &mii.Name, &mii.MiiData)

		friendMii := friends_3ds_types.NewFriendMii()
		friendMii.PID = nex.NewPID(pid)
		friendMii.Mii = mii
		friendMii.ModifiedAt = changedTime

		friendMiis = append(friendMiis, friendMii)
	}

	return friendMiis, nil
}
