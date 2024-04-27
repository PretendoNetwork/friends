package database_3ds

import (
	"github.com/PretendoNetwork/friends/database"
	"github.com/PretendoNetwork/nex-go/v2/types"
	friends_3ds_types "github.com/PretendoNetwork/nex-protocols-go/v2/friends-3ds/types"
)

// UpdateUserMii updates a user's mii
func UpdateUserMii(pid uint32, mii *friends_3ds_types.Mii) error {
	_, err := database.Manager.Exec(`
		INSERT INTO "3ds".user_data (pid, mii_name, mii_data, mii_changed)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (pid)
		DO UPDATE SET 
		mii_name = $2,
		mii_data = $3,
		mii_changed = $4`, pid, mii.Name.Value, mii.MiiData.Value, types.NewDateTime(0).Now().Value())

	if err != nil {
		return err
	}

	return nil
}
