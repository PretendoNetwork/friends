package database_3ds

import (
	"github.com/PretendoNetwork/friends/database"
	"github.com/PretendoNetwork/nex-go/v2/types"
	friends_3ds_types "github.com/PretendoNetwork/nex-protocols-go/v2/friends-3ds/types"
	"github.com/lib/pq"
)

// GetFriendMiis returns the Mii of all friends
func GetFriendMiis(pids []uint32) (*types.List[*friends_3ds_types.FriendMii], error) {
	friendMiis := types.NewList[*friends_3ds_types.FriendMii]()
	friendMiis.Type = friends_3ds_types.NewFriendMii()

	rows, err := database.Postgres.Query(`
	SELECT pid, mii_name, mii_data FROM "3ds".user_data WHERE pid=ANY($1::int[])`, pq.Array(pids))
	if err != nil {
		return friendMiis, err
	}

	changedTime := types.NewDateTime(0).Now()

	for rows.Next() {
		var pid uint32
		var miiName string
		var miiData []byte

		err := rows.Scan(&pid, &miiName, &miiData)
		if err != nil {
			return friendMiis, err
		}

		mii := friends_3ds_types.NewMii()
		mii.Name = types.NewString(miiName)
		mii.Unknown2 = types.NewPrimitiveBool(false)
		mii.Unknown3 = types.NewPrimitiveU8(0)
		mii.MiiData = types.NewBuffer(miiData)

		friendMii := friends_3ds_types.NewFriendMii()
		friendMii.PID = types.NewPID(uint64(pid))
		friendMii.Mii = mii
		friendMii.ModifiedAt = changedTime

		friendMiis.Append(friendMii)
	}

	return friendMiis, nil
}
