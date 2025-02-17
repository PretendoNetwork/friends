package database_wiiu

import (
	"github.com/PretendoNetwork/friends/database"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/v2/friends-wiiu/types"
)

// UpdateUserPrincipalBasicInfo updates the user's basic info
func UpdateUserPrincipalBasicInfo(pid uint32, principalBasicInfo friends_wiiu_types.PrincipalBasicInfo) error {
	_, err := database.Manager.Exec(
		`INSERT INTO wiiu.principal_basic_info (pid, username, unknown)
		VALUES ($1, $2, $3)
		ON CONFLICT (pid)
		DO UPDATE
		SET username = $2, unknown = $3`,
		pid,
		string(principalBasicInfo.NNID),
		uint8(principalBasicInfo.Unknown),
	)
	if err != nil {
		return err
	}

	return UpdateUserMii(pid, principalBasicInfo.Mii)
}
