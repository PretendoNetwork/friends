package database_wiiu

import (
	"github.com/PretendoNetwork/friends-secure/database"
	"github.com/PretendoNetwork/friends-secure/globals"
	"github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

// Get a users blacklist
func GetUserBlockList(pid uint32) []*nexproto.BlacklistedPrincipal {
	blockList := make([]*nexproto.BlacklistedPrincipal, 0)

	rows, err := database.Postgres.Query(`SELECT blocked_pid, title_id, title_version, date FROM wiiu.blocks WHERE blocker_pid=$1`, pid)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return blockList
	}

	for rows.Next() {
		var pid uint32
		var titleId uint64
		var titleVersion uint16
		var date *nex.DateTime
		rows.Scan(&pid, &titleId, &titleVersion, &date)

		blacklistPrincipal := nexproto.NewBlacklistedPrincipal()

		blacklistPrincipal.PrincipalBasicInfo = GetUserInfoByPID(pid)

		blacklistPrincipal.GameKey = nexproto.NewGameKey()
		blacklistPrincipal.GameKey.TitleID = titleId
		blacklistPrincipal.GameKey.TitleVersion = titleVersion
		blacklistPrincipal.BlackListedSince = date

		blockList = append(blockList, blacklistPrincipal)
	}

	return blockList
}
