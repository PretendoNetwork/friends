package database_wiiu

import (
	"database/sql"

	"github.com/PretendoNetwork/friends/database"
	"github.com/PretendoNetwork/friends/utility"
	"github.com/PretendoNetwork/nex-go/v2/types"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/v2/friends-wiiu/types"
)

// GetUserBlockList returns a user's blacklist
func GetUserBlockList(pid uint32) (*types.List[*friends_wiiu_types.BlacklistedPrincipal], error) {
	blockList := types.NewList[*friends_wiiu_types.BlacklistedPrincipal]()
	blockList.Type = friends_wiiu_types.NewBlacklistedPrincipal()

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
		var titleID uint64
		var titleVersion uint16
		var date uint64

		err := rows.Scan(&pid, &titleID, &titleVersion, &date)
		if err != nil {
			return nil, err
		}

		userInfo, err := utility.GetUserInfoByPID(pid)
		if err != nil {
			return nil, err
		}

		blacklistPrincipal := friends_wiiu_types.NewBlacklistedPrincipal()

		blacklistPrincipal.PrincipalBasicInfo = userInfo

		blacklistPrincipal.GameKey = friends_wiiu_types.NewGameKey()
		blacklistPrincipal.GameKey.TitleID = types.NewPrimitiveU64(titleID)
		blacklistPrincipal.GameKey.TitleVersion = types.NewPrimitiveU16(titleVersion)
		blacklistPrincipal.BlackListedSince = types.NewDateTime(date)

		blockList.Append(blacklistPrincipal)
	}

	return blockList, nil
}
