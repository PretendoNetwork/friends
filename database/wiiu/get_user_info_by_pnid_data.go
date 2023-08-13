package database_wiiu

import (
	"encoding/base64"

	pb "github.com/PretendoNetwork/grpc-go/account"
	"github.com/PretendoNetwork/nex-go"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/friends-wiiu/types"
)

// GetUserInfoByPNIDData converts the account's PNID data into user info for friends
func GetUserInfoByPNIDData(userData *pb.GetUserDataResponse) (*friends_wiiu_types.PrincipalBasicInfo, error) {
	info := friends_wiiu_types.NewPrincipalBasicInfo()

	info.PID = userData.Pid
	info.NNID = userData.Username
	info.Mii = friends_wiiu_types.NewMiiV2()
	info.Unknown = 2

	encodedMiiData := userData.Mii.Data
	decodedMiiData, err := base64.StdEncoding.DecodeString(encodedMiiData)
	if err != nil {
		return nil, err
	}

	info.Mii.Name = userData.Mii.Name
	info.Mii.Unknown1 = 0
	info.Mii.Unknown2 = 0
	info.Mii.MiiData = decodedMiiData
	info.Mii.Datetime = nex.NewDateTime(0)

	return info, nil
}
