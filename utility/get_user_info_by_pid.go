package utility

import (
	"encoding/base64"

	"github.com/PretendoNetwork/nex-go/v2/types"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/v2/friends-wiiu/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/PretendoNetwork/friends/database"
	"github.com/PretendoNetwork/friends/globals"
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

	info.PID = types.NewPID(uint64(userData.Pid))
	info.NNID = types.NewString(userData.Username)
	info.Mii = friends_wiiu_types.NewMiiV2()
	info.Unknown = types.NewPrimitiveU8(2)

	encodedMiiData := userData.Mii.Data
	decodedMiiData, err := base64.StdEncoding.DecodeString(encodedMiiData)
	if err != nil {
		return nil, err
	}

	info.Mii.Name = types.NewString(userData.Mii.Name)
	info.Mii.Unknown1 = types.NewPrimitiveU8(0)
	info.Mii.Unknown2 = types.NewPrimitiveU8(0)
	info.Mii.MiiData = types.NewBuffer(decodedMiiData)
	info.Mii.Datetime = types.NewDateTime(0)

	return info, nil
}
