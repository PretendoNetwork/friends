package database_wiiu

import (
	"database/sql"

	"github.com/PretendoNetwork/friends/database"
	"github.com/PretendoNetwork/nex-go/v2/types"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/v2/friends-wiiu/types"
)

// GetUserBlockList returns a user's blacklist
func GetUserBlockList(pid uint32) (*types.List[*friends_wiiu_types.BlacklistedPrincipal], error) {
	blockList := types.NewList[*friends_wiiu_types.BlacklistedPrincipal]()
	blockList.Type = friends_wiiu_types.NewBlacklistedPrincipal()

	rows, err := database.Manager.Query(`
	SELECT
		b.blocked_pid, b.title_id, b.title_version, b.date,
		bi.username, bi.unknown,
		mii.name, mii.unknown1, mii.unknown2, mii.data, mii.unknown_datetime
	FROM wiiu.blocks AS b
	INNER JOIN wiiu.principal_basic_info AS bi ON bi.pid = b.blocked_pid
	INNER JOIN wiiu.mii AS mii ON mii.pid = b.blocked_pid
	WHERE blocker_pid=$1
	`, pid)
	if err != nil {
		if err == sql.ErrNoRows {
			return blockList, database.ErrBlacklistNotFound
		} else {
			return blockList, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		var blockedPID uint32
		var titleID uint64
		var titleVersion uint16
		var date uint64

		var blockedNNID string
		var unknown uint8

		var miiName string
		var miiUnknown1 uint8
		var miiUnknown2 uint8
		var miiData []byte
		var miiDatetime uint64

		err := rows.Scan(&blockedPID, &titleID, &titleVersion, &date, &blockedNNID, &unknown, &miiName, &miiUnknown1, &miiUnknown2, &miiData, &miiDatetime)
		if err != nil {
			return nil, err
		}

		mii := friends_wiiu_types.NewMiiV2()
		mii.Name = types.NewString(miiName)
		mii.Unknown1 = types.NewPrimitiveU8(miiUnknown1)
		mii.Unknown2 = types.NewPrimitiveU8(miiUnknown2)
		mii.MiiData = types.NewBuffer(miiData)
		mii.Datetime = types.NewDateTime(miiDatetime)

		principalBasicInfo := friends_wiiu_types.NewPrincipalBasicInfo()
		principalBasicInfo.PID = types.NewPID(uint64(blockedPID))
		principalBasicInfo.NNID = types.NewString(blockedNNID)
		principalBasicInfo.Unknown = types.NewPrimitiveU8(unknown)
		principalBasicInfo.Mii = mii

		blacklistPrincipal := friends_wiiu_types.NewBlacklistedPrincipal()
		blacklistPrincipal.PrincipalBasicInfo = principalBasicInfo
		blacklistPrincipal.GameKey = friends_wiiu_types.NewGameKey()
		blacklistPrincipal.GameKey.TitleID = types.NewPrimitiveU64(titleID)
		blacklistPrincipal.GameKey.TitleVersion = types.NewPrimitiveU16(titleVersion)
		blacklistPrincipal.BlackListedSince = types.NewDateTime(date)

		blockList.Append(blacklistPrincipal)
	}

	return blockList, nil
}
