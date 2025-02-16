package database_wiiu

import (
	"database/sql"
	"time"

	"github.com/PretendoNetwork/friends/database"
	"github.com/PretendoNetwork/nex-go/v2/types"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/v2/friends-wiiu/types"
)

// GetUserFriendRequestsIn returns the friend requests received by a user
func GetUserFriendRequestsIn(pid uint32) (types.List[friends_wiiu_types.FriendRequest], error) {
	friendRequests := types.NewList[friends_wiiu_types.FriendRequest]()

	rows, err := database.Manager.Query(`
	SELECT
		fr.id, fr.sender_pid, fr.sent_on, fr.expires_on, fr.message, fr.received,
		bi.username, bi.unknown,
		mii.name, mii.unknown1, mii.unknown2, mii.data, mii.unknown_datetime
	FROM wiiu.friend_requests AS fr
	INNER JOIN wiiu.principal_basic_info AS bi ON bi.pid = fr.sender_pid
	INNER JOIN wiiu.mii AS mii ON mii.pid = fr.sender_pid
	WHERE recipient_pid=$1 AND accepted=false AND denied=false
	LIMIT 100
	`, pid)
	if err != nil {
		if err == sql.ErrNoRows {
			return friendRequests, database.ErrEmptyList
		} else {
			return friendRequests, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		var id uint64
		var senderPID uint32
		var sentOn uint64
		var expiresOn uint64
		var message string
		var received bool

		var senderNNID string
		var unknown uint8

		var miiName string
		var miiUnknown1 uint8
		var miiUnknown2 uint8
		var miiData []byte
		var miiDatetime uint64

		err := rows.Scan(&id, &senderPID, &sentOn, &expiresOn, &message, &received, &senderNNID, &unknown, &miiName, &miiUnknown1, &miiUnknown2, &miiData, &miiDatetime)
		if err != nil {
			return friendRequests, err
		}

		mii := friends_wiiu_types.NewMiiV2()
		mii.Name = types.NewString(miiName)
		mii.Unknown1 = types.NewUInt8(miiUnknown1)
		mii.Unknown2 = types.NewUInt8(miiUnknown2)
		mii.MiiData = types.NewBuffer(miiData)
		mii.Datetime = types.NewDateTime(miiDatetime)

		principalBasicInfo := friends_wiiu_types.NewPrincipalBasicInfo()
		principalBasicInfo.PID = types.NewPID(uint64(senderPID))
		principalBasicInfo.NNID = types.NewString(senderNNID)
		principalBasicInfo.Unknown = types.NewUInt8(unknown)
		principalBasicInfo.Mii = mii

		friendRequest := friends_wiiu_types.NewFriendRequest()
		friendRequest.PrincipalInfo = principalBasicInfo
		friendRequest.Message = friends_wiiu_types.NewFriendRequestMessage()
		friendRequest.Message.FriendRequestID = types.NewUInt64(id)
		friendRequest.Message.Received = types.NewBool(received)
		friendRequest.Message.Unknown2 = types.NewUInt8(1)
		friendRequest.Message.Message = types.NewString(message)
		friendRequest.Message.Unknown3 = types.NewUInt8(0)
		friendRequest.Message.Unknown4 = types.NewString("")
		friendRequest.Message.GameKey = friends_wiiu_types.NewGameKey()
		friendRequest.Message.GameKey.TitleID = types.NewUInt64(0)
		friendRequest.Message.GameKey.TitleVersion = types.NewUInt16(0)
		friendRequest.Message.Unknown5 = types.NewDateTime(134222053376) // * idk what this value means but its always this
		friendRequest.Message.ExpiresOn = types.NewDateTime(expiresOn)
		friendRequest.SentOn = types.NewDateTime(sentOn)

		// * Filter out expired requests
		if friendRequest.Message.ExpiresOn.Standard().After(time.Now()) {
			friendRequests = append(friendRequests, friendRequest)
		}
	}

	return friendRequests, nil
}
