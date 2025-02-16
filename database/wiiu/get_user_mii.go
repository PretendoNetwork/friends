package database_wiiu

import (
	"database/sql"

	"github.com/PretendoNetwork/friends/database"
	"github.com/PretendoNetwork/nex-go/v2/types"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/v2/friends-wiiu/types"
)

// GetUserMii returns the users Mii
func GetUserMii(pid uint32) (friends_wiiu_types.MiiV2, error) {
	mii := friends_wiiu_types.NewMiiV2()

	var name string
	var unknown1 uint8
	var unknown2 uint8
	var data []byte
	var datetime uint64

	row, err := database.Manager.QueryRow(`SELECT name, unknown1, unknown2, data, unknown_datetime FROM wiiu.mii WHERE pid=$1`, pid)
	if err != nil {
		return mii, err
	}

	err = row.Scan(&name, &unknown1, &unknown2, &data, &datetime)
	if err != nil {
		if err == sql.ErrNoRows {
			return mii, database.ErrPIDNotFound
		} else {
			return mii, err
		}
	}

	mii.Name = types.NewString(name)
	mii.Unknown1 = types.NewUInt8(unknown1)
	mii.Unknown2 = types.NewUInt8(unknown2)
	mii.MiiData = types.NewBuffer(data)
	mii.Datetime = types.NewDateTime(datetime)

	return mii, nil
}
