package database_wiiu

import (
	"github.com/PretendoNetwork/friends/database"
	"github.com/PretendoNetwork/nex-go/v2/types"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/v2/friends-wiiu/types"
)

// UpdateNetworkAccountInfo updates a user's network account information
func UpdateNetworkAccountInfo(pid uint32, nnaInfo friends_wiiu_types.NNAInfo, birthday types.DateTime) error {
	_, err := database.Manager.Exec(
		`INSERT INTO wiiu.network_account_info (pid, unknown1, unknown2, birthday)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (pid)
		DO UPDATE
		SET unknown1 = $2, unknown2 = $3, birthday = $4`,
		pid,
		uint8(nnaInfo.Unknown1),
		uint8(nnaInfo.Unknown2),
		uint64(birthday),
	)
	if err != nil {
		return err
	}

	return UpdateUserPrincipalBasicInfo(pid, nnaInfo.PrincipalBasicInfo)
}
