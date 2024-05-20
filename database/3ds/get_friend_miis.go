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

	rows, err := database.Manager.Query(`
	SELECT pid, mii_name, mii_profanity, mii_character_set, mii_data, mii_changed FROM "3ds".user_data WHERE pid=ANY($1::int[])`, pq.Array(pids))
	if err != nil {
		return friendMiis, err
	}
	defer rows.Close()

	for rows.Next() {
		var pid uint32
		var miiName string
		var miiProfanity bool
		var miiCharacterSet uint8
		var miiData []byte
		var changedTime uint64

		err := rows.Scan(&pid, &miiName, &miiProfanity, &miiCharacterSet, &miiData, &changedTime)
		if err != nil {
			return friendMiis, err
		}

		mii := friends_3ds_types.NewMii()
		mii.Name = types.NewString(miiName)
		mii.ProfanityFlag = types.NewPrimitiveBool(miiProfanity)
		mii.CharacterSet = types.NewPrimitiveU8(miiCharacterSet)
		mii.MiiData = types.NewBuffer(miiData)

		friendMii := friends_3ds_types.NewFriendMii()
		friendMii.PID = types.NewPID(uint64(pid))
		friendMii.Mii = mii
		friendMii.ModifiedAt = types.NewDateTime(changedTime)

		friendMiis.Append(friendMii)
	}

	return friendMiis, nil
}
