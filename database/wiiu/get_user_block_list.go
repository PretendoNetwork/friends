package database_wiiu

import (
	"github.com/PretendoNetwork/friends-secure/database"
	"github.com/PretendoNetwork/friends-secure/globals"
	"github.com/PretendoNetwork/nex-go"
	friends_wiiu "github.com/PretendoNetwork/nex-protocols-go/friends/wiiu"
)

// Get a users blacklist
func GetUserBlockList(pid uint32) []*friends_wiiu.BlacklistedPrincipal {
	blockList := make([]*friends_wiiu.BlacklistedPrincipal, 0)

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

		blacklistPrincipal := friends_wiiu.NewBlacklistedPrincipal()

		blacklistPrincipal.PrincipalBasicInfo = GetUserInfoByPID(pid)

		blacklistPrincipal.GameKey = friends_wiiu.NewGameKey()
		blacklistPrincipal.GameKey.TitleID = titleId
		blacklistPrincipal.GameKey.TitleVersion = titleVersion
		blacklistPrincipal.BlackListedSince = date

		blockList = append(blockList, blacklistPrincipal)
	}

	return blockList
}
