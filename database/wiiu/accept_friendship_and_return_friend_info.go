package database_wiiu

import (
	"encoding/base64"
	"time"

	"github.com/PretendoNetwork/friends-secure/database"
	"github.com/PretendoNetwork/friends-secure/globals"
	"github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
	"github.com/gocql/gocql"
	"go.mongodb.org/mongo-driver/bson"
)

func AcceptFriendshipAndReturnFriendInfo(friendRequestID uint64) *nexproto.FriendInfo {
	var senderPID uint32
	var recipientPID uint32

	err := database.Postgres.QueryRow(`SELECT sender_pid, recipient_pid FROM wiiu.friend_requests WHERE id=$1`, friendRequestID).Scan(&senderPID, &recipientPID)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return nil
	}

	acceptedTime := nex.NewDateTime(0)
	acceptedTime.FromTimestamp(time.Now())

	// Friendships are two-way relationships, not just one link between 2 entities
	// "A" has friend "B" and "B" has friend "A", so store both relationships

	_, err = database.Postgres.Exec(`INSERT INTO wiiu.friendships (user1_pid, user2_pid, date) VALUES ($1, $2, $3)`, senderPID, recipientPID, acceptedTime.Value())
	if err != nil {
		globals.Logger.Critical(err.Error())
		return nil
	}

	_, err = database.Postgres.Exec(`INSERT INTO wiiu.friendships (user1_pid, user2_pid, date) VALUES ($1, $2, $3)`, recipientPID, senderPID, acceptedTime.Value())
	if err != nil {
		globals.Logger.Critical(err.Error())
		return nil
	}

	SetFriendRequestAccepted(friendRequestID)

	friendInfo := nexproto.NewFriendInfo()
	connectedUser := globals.ConnectedUsers[senderPID]
	var lastOnline *nex.DateTime

	if connectedUser != nil {
		// Online
		friendInfo.NNAInfo = connectedUser.NNAInfo
		friendInfo.Presence = connectedUser.Presence

		lastOnline = nex.NewDateTime(0)
		lastOnline.FromTimestamp(time.Now())
	} else {
		// Offline
		senderUserInforation := GetUserInfoByPID(senderPID)
		encodedMiiData := senderUserInforation["mii"].(bson.M)["data"].(string)
		decodedMiiData, _ := base64.StdEncoding.DecodeString(encodedMiiData)

		friendInfo.NNAInfo = nexproto.NewNNAInfo()
		friendInfo.NNAInfo.PrincipalBasicInfo = nexproto.NewPrincipalBasicInfo()
		friendInfo.NNAInfo.PrincipalBasicInfo.PID = senderPID
		friendInfo.NNAInfo.PrincipalBasicInfo.NNID = senderUserInforation["username"].(string)
		friendInfo.NNAInfo.PrincipalBasicInfo.Mii = nexproto.NewMiiV2()
		friendInfo.NNAInfo.PrincipalBasicInfo.Mii.Name = senderUserInforation["mii"].(bson.M)["name"].(string)
		friendInfo.NNAInfo.PrincipalBasicInfo.Mii.Unknown1 = 0
		friendInfo.NNAInfo.PrincipalBasicInfo.Mii.Unknown2 = 0
		friendInfo.NNAInfo.PrincipalBasicInfo.Mii.Data = decodedMiiData
		friendInfo.NNAInfo.PrincipalBasicInfo.Mii.Datetime = nex.NewDateTime(0)
		friendInfo.NNAInfo.PrincipalBasicInfo.Unknown = 0
		friendInfo.NNAInfo.Unknown1 = 0
		friendInfo.NNAInfo.Unknown2 = 0

		friendInfo.Presence = nexproto.NewNintendoPresenceV2()
		friendInfo.Presence.ChangedFlags = 0
		friendInfo.Presence.Online = false
		friendInfo.Presence.GameKey = nexproto.NewGameKey()
		friendInfo.Presence.GameKey.TitleID = 0
		friendInfo.Presence.GameKey.TitleVersion = 0
		friendInfo.Presence.Unknown1 = 0
		friendInfo.Presence.Message = ""
		friendInfo.Presence.Unknown2 = 0
		friendInfo.Presence.Unknown3 = 0
		friendInfo.Presence.GameServerID = 0
		friendInfo.Presence.Unknown4 = 0
		friendInfo.Presence.PID = senderPID
		friendInfo.Presence.GatheringID = 0
		friendInfo.Presence.ApplicationData = []byte{0x00}
		friendInfo.Presence.Unknown5 = 0
		friendInfo.Presence.Unknown6 = 0
		friendInfo.Presence.Unknown7 = 0

		var lastOnlineTime uint64
		err := database.Postgres.QueryRow(`SELECT last_online FROM wiiu.user_data WHERE pid=$1`, senderPID).Scan(&lastOnlineTime)
		if err != nil {
			lastOnlineTime = nex.NewDateTime(0).Now()

			if err == gocql.ErrNotFound {
				globals.Logger.Error(err.Error())
			} else {
				globals.Logger.Critical(err.Error())
			}
		}

		lastOnline = nex.NewDateTime(lastOnlineTime) // TODO: Change this
	}

	friendInfo.Status = GetUserComment(senderPID)
	friendInfo.BecameFriend = acceptedTime
	friendInfo.LastOnline = lastOnline // TODO: Change this
	friendInfo.Unknown = 0

	return friendInfo
}
