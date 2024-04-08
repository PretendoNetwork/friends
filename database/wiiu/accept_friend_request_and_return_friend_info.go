package database_wiiu

import (
	"database/sql"

	"github.com/PretendoNetwork/friends/database"
	"github.com/PretendoNetwork/friends/globals"
	"github.com/PretendoNetwork/friends/utility"
	"github.com/PretendoNetwork/nex-go/v2/types"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/v2/friends-wiiu/types"
)

// AcceptFriendRequestAndReturnFriendInfo accepts the given friend reuqest and returns the friend's information
func AcceptFriendRequestAndReturnFriendInfo(friendRequestID uint64) (*friends_wiiu_types.FriendInfo, error) {
	var senderPID uint32
	var recipientPID uint32

	err := database.Postgres.QueryRow(`SELECT sender_pid, recipient_pid FROM wiiu.friend_requests WHERE id=$1`, friendRequestID).Scan(&senderPID, &recipientPID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, database.ErrFriendRequestNotFound
		} else {
			return nil, err
		}
	}

	acceptedTime := types.NewDateTime(0).Now()

	// * Friendships are two-way relationships, not just one link between 2 entities
	// * "A" has friend "B" and "B" has friend "A", so store both relationships

	// * If were friends before, just activate the status again

	_, err = database.Postgres.Exec(`
		INSERT INTO wiiu.friendships (user1_pid, user2_pid, date, active)
		VALUES ($1, $2, $3, true)
		ON CONFLICT (user1_pid, user2_pid)
		DO UPDATE SET
		date = $3,
		active = true`, senderPID, recipientPID, acceptedTime.Value())
	if err != nil {
		return nil, err
	}

	_, err = database.Postgres.Exec(`
		INSERT INTO wiiu.friendships (user1_pid, user2_pid, date, active)
		VALUES ($1, $2, $3, true)
		ON CONFLICT (user1_pid, user2_pid)
		DO UPDATE SET
		date = $3,
		active = true`, recipientPID, senderPID, acceptedTime.Value())
	if err != nil {
		return nil, err
	}

	err = SetFriendRequestAccepted(friendRequestID)
	if err != nil {
		return nil, err
	}

	friendInfo := friends_wiiu_types.NewFriendInfo()
	connectedUser := globals.ConnectedUsers[senderPID]
	var lastOnline *types.DateTime

	if connectedUser != nil {
		// * Online
		friendInfo.NNAInfo = connectedUser.NNAInfo
		friendInfo.Presence = connectedUser.PresenceV2

		lastOnline = types.NewDateTime(0).Now()
	} else {
		// * Offline
		userInfo, err := utility.GetUserInfoByPID(senderPID)
		if err != nil {
			return nil, err
		}

		friendInfo.NNAInfo = friends_wiiu_types.NewNNAInfo()

		friendInfo.NNAInfo.PrincipalBasicInfo = userInfo
		friendInfo.NNAInfo.Unknown1 = types.NewPrimitiveU8(0)
		friendInfo.NNAInfo.Unknown2 = types.NewPrimitiveU8(0)

		friendInfo.Presence = friends_wiiu_types.NewNintendoPresenceV2()
		friendInfo.Presence.ChangedFlags = types.NewPrimitiveU32(0)
		friendInfo.Presence.Online = types.NewPrimitiveBool(false)
		friendInfo.Presence.GameKey = friends_wiiu_types.NewGameKey()
		friendInfo.Presence.GameKey.TitleID = types.NewPrimitiveU64(0)
		friendInfo.Presence.GameKey.TitleVersion = types.NewPrimitiveU16(0)
		friendInfo.Presence.Unknown1 = types.NewPrimitiveU8(0)
		friendInfo.Presence.Message = types.NewString("")
		friendInfo.Presence.Unknown2 = types.NewPrimitiveU32(0)
		friendInfo.Presence.Unknown3 = types.NewPrimitiveU8(0)
		friendInfo.Presence.GameServerID = types.NewPrimitiveU32(0)
		friendInfo.Presence.Unknown4 = types.NewPrimitiveU32(0)
		friendInfo.Presence.PID = types.NewPID(uint64(senderPID))
		friendInfo.Presence.GatheringID = types.NewPrimitiveU32(0)
		friendInfo.Presence.ApplicationData = types.NewBuffer([]byte{})
		friendInfo.Presence.Unknown5 = types.NewPrimitiveU8(0)
		friendInfo.Presence.Unknown6 = types.NewPrimitiveU8(0)
		friendInfo.Presence.Unknown7 = types.NewPrimitiveU8(0)

		var lastOnlineTime uint64
		err = database.Postgres.QueryRow(`SELECT last_online FROM wiiu.user_data WHERE pid=$1`, senderPID).Scan(&lastOnlineTime)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, database.ErrPIDNotFound
			} else {
				return nil, err
			}
		}

		lastOnline = types.NewDateTime(lastOnlineTime) // TODO - Change this
	}

	status, err := GetUserComment(senderPID)
	if err != nil {
		return nil, err
	}

	friendInfo.Status = status
	friendInfo.BecameFriend = acceptedTime
	friendInfo.LastOnline = lastOnline // TODO - Change this
	friendInfo.Unknown = types.NewPrimitiveU64(0)

	return friendInfo, nil
}
