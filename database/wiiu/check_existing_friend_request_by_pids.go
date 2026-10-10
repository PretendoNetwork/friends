package database_wiiu

import (
	"database/sql"

	"github.com/PretendoNetwork/friends/database"
	"github.com/PretendoNetwork/nex-go/v2/types"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/v2/friends-wiiu/types"
)

// CheckExistingFriendRequestByPID checks for an existing friend request from a PID pair
func CheckExistingFriendRequestByPIDs(targetPID uint32, userPID uint32) (uint64, friends_wiiu_types.FriendRequestMessage, error) {
	friendRequestMessage := friends_wiiu_types.NewFriendRequestMessage()

	row, err := database.Manager.QueryRow(`
	SELECT
		fr.id, fr.message, fr.provisional
	FROM wiiu.friend_requests AS fr
	WHERE sender_pid=$1 AND recipient_pid=$2 AND accepted=false AND denied=false
	`, targetPID, userPID)
	if err != nil {
		return 0, friendRequestMessage, err
	}

	var id uint64
	var message string
	var provisional bool

	err = row.Scan(&id, &message, &provisional)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, friendRequestMessage, database.ErrFriendRequestNotFound
		} else {
			return 0, friendRequestMessage, err
		}
	}

	// * All hardcoded except Message. Matching NN

	friendRequestMessage.FriendRequestID = types.NewUInt64(0)
	friendRequestMessage.Received = types.NewBool(false)
	friendRequestMessage.Unknown2 = types.NewUInt8(1)
	friendRequestMessage.Message = types.NewString(message)
	friendRequestMessage.Unknown3 = types.NewUInt8(0)
	friendRequestMessage.Unknown4 = types.NewString("")
	friendRequestMessage.GameKey = friends_wiiu_types.NewGameKey()
	friendRequestMessage.GameKey.TitleID = types.NewUInt64(0)
	friendRequestMessage.GameKey.TitleVersion = types.NewUInt16(0)
	friendRequestMessage.Unknown5 = types.NewDateTime(0)
	friendRequestMessage.ExpiresOn = types.NewDateTime(0)

	return id, friendRequestMessage, nil
}
