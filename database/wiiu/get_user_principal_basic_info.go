package database_wiiu

import (
	"database/sql"

	"github.com/PretendoNetwork/friends/database"
	"github.com/PretendoNetwork/nex-go/v2/types"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/v2/friends-wiiu/types"
)

// GetUserPrincipalBasicInfo returns the users basic info
func GetUserPrincipalBasicInfo(pid uint32) (friends_wiiu_types.PrincipalBasicInfo, error) {
	principalBasicInfo := friends_wiiu_types.NewPrincipalBasicInfo()

	var nnid string
	var unknown uint8

	row, err := database.Manager.QueryRow(`SELECT username, unknown FROM wiiu.principal_basic_info WHERE pid=$1`, pid)
	if err != nil {
		return principalBasicInfo, err
	}

	err = row.Scan(&nnid, &unknown)
	if err != nil {
		if err == sql.ErrNoRows {
			return principalBasicInfo, database.ErrPIDNotFound
		} else {
			return principalBasicInfo, err
		}
	}

	principalBasicInfo.PID = types.NewPID(uint64(pid))
	principalBasicInfo.NNID = types.NewString(nnid)
	principalBasicInfo.Unknown = types.NewUInt8(unknown)
	principalBasicInfo.Mii, err = GetUserMii(pid)
	if err != nil {
		return principalBasicInfo, err
	}

	return principalBasicInfo, nil
}
