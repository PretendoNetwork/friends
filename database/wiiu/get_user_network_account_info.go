package database_wiiu

import (
	"database/sql"

	"github.com/PretendoNetwork/friends/database"
	"github.com/PretendoNetwork/nex-go/v2/types"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/v2/friends-wiiu/types"
)

// GetUserNetworkAccountInfo returns the users network account info
func GetUserNetworkAccountInfo(pid uint32) (friends_wiiu_types.NNAInfo, error) {
	nnaInfo := friends_wiiu_types.NewNNAInfo()

	var unknown1 uint8
	var unknown2 uint8

	row, err := database.Manager.QueryRow(`SELECT unknown1, unknown2 FROM wiiu.network_account_info WHERE pid=$1`, pid)
	if err != nil {
		return nnaInfo, err
	}

	err = row.Scan(&unknown1, &unknown2)
	if err != nil {
		if err == sql.ErrNoRows {
			return nnaInfo, database.ErrPIDNotFound
		} else {
			return nnaInfo, err
		}
	}

	nnaInfo.Unknown1 = types.NewUInt8(unknown1)
	nnaInfo.Unknown2 = types.NewUInt8(unknown2)
	nnaInfo.PrincipalBasicInfo, err = GetUserPrincipalBasicInfo(pid)
	if err != nil {
		return nnaInfo, err
	}

	return nnaInfo, nil
}
