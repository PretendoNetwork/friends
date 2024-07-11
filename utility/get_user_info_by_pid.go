package utility

import (
	"encoding/base64"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"github.com/PretendoNetwork/nex-go"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/friends-wiiu/types"

	"github.com/CloudnetworkTeam/friends/database"
	"github.com/CloudnetworkTeam/friends/globals"
)

// GetUserInfoByPID returns the user information for a given PID
func GetUserInfoByPID(pid uint32) (*friends_wiiu_types.PrincipalBasicInfo, error) {
	userData, err := globals.GetUserData(pid)
	if err != nil {
		if status.Code(err) == codes.InvalidArgument {
			return nil, database.ErrPIDNotFound
		} else {
			return nil, err
		}
	}

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
