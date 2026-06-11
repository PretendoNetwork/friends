package database_3ds

import (
	"database/sql"

	"github.com/PretendoNetwork/friends/database"
	"github.com/PretendoNetwork/nex-go/v2/types"
	friends_3ds_types "github.com/PretendoNetwork/nex-protocols-go/v2/friends-3ds/types"
)

// GetMii returns the Mii of a specified user
func GetMii(pid types.PID) (friends_3ds_types.FriendMii, error) {
	friendMii := friends_3ds_types.NewFriendMii()

	rows, err := database.Manager.QueryRow(`
	SELECT mii_name, mii_profanity, mii_character_set, mii_data, mii_changed FROM "3ds".user_data WHERE pid=$1`, pid)
	if err != nil {
		return friendMii, err
	}

	var miiName string
	var miiProfanity bool
	var miiCharacterSet uint8
	var miiData []byte
	var changedTime uint64

	err = rows.Scan(&miiName, &miiProfanity, &miiCharacterSet, &miiData, &changedTime)
	if err != nil {
		if err == sql.ErrNoRows {
			return friendMii, database.ErrPIDNotFound
		} else {
			return friendMii, err
		}
	}

	friendMii.PID = types.NewPID(uint64(pid))
	friendMii.Mii.Name = types.NewString(miiName)
	friendMii.Mii.ProfanityFlag = types.NewBool(miiProfanity)
	friendMii.Mii.CharacterSet = types.NewUInt8(miiCharacterSet)
	friendMii.Mii.MiiData = types.NewBuffer(miiData)
	friendMii.ModifiedAt = types.NewDateTime(changedTime)

	return friendMii, nil
}
