package database_wiiu

import (
	"database/sql"

	"github.com/PretendoNetwork/friends/database"
	"github.com/PretendoNetwork/nex-go/v2/types"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/v2/friends-wiiu/types"
)

// GetUserPrincipalBasicInfoByNNID returns the users basic info from the user's NNID
func GetUserPrincipalBasicInfoByNNID(nnid types.String) (friends_wiiu_types.PrincipalBasicInfo, error) {
	principalBasicInfo := friends_wiiu_types.NewPrincipalBasicInfo()

	var pid types.PID
	var unknown uint8

	row, err := database.Manager.QueryRow(`SELECT pid, unknown FROM wiiu.principal_basic_info WHERE username=$1`, nnid)
	if err != nil {
		return principalBasicInfo, err
	}

	err = row.Scan(&pid, &unknown)
	if err != nil {
		if err == sql.ErrNoRows {
			return principalBasicInfo, database.ErrPIDNotFound
		} else {
			return principalBasicInfo, err
		}
	}

	principalBasicInfo.PID = types.NewPID(uint64(pid))
	principalBasicInfo.NNID = nnid
	principalBasicInfo.Unknown = types.NewUInt8(unknown)
	principalBasicInfo.Mii, err = GetUserMii(uint32(pid))
	if err != nil {
		return principalBasicInfo, err
	}

	return principalBasicInfo, nil
}
