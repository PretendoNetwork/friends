package database_wiiu

import (
	"database/sql"

	"github.com/PretendoNetwork/friends/database"
	"github.com/PretendoNetwork/friends/globals"
	"github.com/PretendoNetwork/nex-go"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/friends-wiiu/types"
)

// GetUserBlockList returns a user's blacklist
func GetUserBlockList(pid uint32) ([]*friends_wiiu_types.BlacklistedPrincipal, error) {
	blockList := make([]*friends_wiiu_types.BlacklistedPrincipal, 0)

	rows, err := database.Postgres.Query(`SELECT blocked_pid, title_id, title_version, date FROM wiiu.blocks WHERE blocker_pid=$1`, pid)
	if err != nil {
		if err == sql.ErrNoRows {
			return blockList, database.ErrBlacklistNotFound
		} else {
			return blockList, err
		}
	}

	for rows.Next() {
		var pid uint32
		var titleId uint64
		var titleVersion uint16
		var date *nex.DateTime
		rows.Scan(&pid, &titleId, &titleVersion, &date)

		userData, err := globals.GetUserData(pid)
		if err != nil {
			globals.Logger.Critical(err.Error())
			continue
		}

		userInfo, err := GetUserInfoByPNIDData(userData)
		if err != nil {
			globals.Logger.Critical(err.Error())
			continue
		}

		blacklistPrincipal := friends_wiiu_types.NewBlacklistedPrincipal()

		blacklistPrincipal.PrincipalBasicInfo = userInfo

		blacklistPrincipal.GameKey = friends_wiiu_types.NewGameKey()
		blacklistPrincipal.GameKey.TitleID = titleId
		blacklistPrincipal.GameKey.TitleVersion = titleVersion
		blacklistPrincipal.BlackListedSince = date

		blockList = append(blockList, blacklistPrincipal)
	}

	return blockList, nil
}
