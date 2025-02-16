package database_wiiu

import (
	"github.com/PretendoNetwork/friends/database"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/v2/friends-wiiu/types"
)

// UpdateUserMii updates the user's Mii
func UpdateUserMii(pid uint32, mii friends_wiiu_types.MiiV2) error {
	_, err := database.Manager.Exec(`
		INSERT INTO wiiu.mii (pid, name, unknown1, unknown2, data, unknown_datetime)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (pid)
		DO UPDATE
		SET
		name = $2,
		unknown1 = $3,
		unknown2 = $4,
		data = $5,
		unknown_datetime = $6`,
		pid,
		string(mii.Name),
		uint8(mii.Unknown1),
		uint8(mii.Unknown2),
		[]byte(mii.MiiData),
		uint64(mii.Datetime),
	)
	if err != nil {
		return err
	}

	return nil
}
