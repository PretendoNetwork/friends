package database_wiiu

import (
	"encoding/base64"

	"github.com/PretendoNetwork/friends/globals"
	"github.com/PretendoNetwork/nex-go"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/friends-wiiu/types"
)

func GetUserInfoByPID(pid uint32) *friends_wiiu_types.PrincipalBasicInfo {
	info := friends_wiiu_types.NewPrincipalBasicInfo()

	userData, err := globals.GetUserData(pid)

	if err != nil {
		globals.Logger.Critical(err.Error())

		return info
	}

	info.PID = pid
	info.NNID = userData.Username
	info.Mii = friends_wiiu_types.NewMiiV2()
	info.Unknown = 2

	encodedMiiData := userData.Mii.Data
	decodedMiiData, _ := base64.StdEncoding.DecodeString(encodedMiiData)

	info.Mii.Name = userData.Mii.Name
	info.Mii.Unknown1 = 0
	info.Mii.Unknown2 = 0
	info.Mii.MiiData = decodedMiiData
	info.Mii.Datetime = nex.NewDateTime(0)

	return info
}
