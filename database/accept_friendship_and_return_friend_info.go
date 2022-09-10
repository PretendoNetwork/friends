package database

import (
	"encoding/base64"
	"math/rand"
	"time"

	"github.com/PretendoNetwork/friends-secure/globals"
	"github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
	"github.com/gocql/gocql"
	"go.mongodb.org/mongo-driver/bson"
)

func AcceptFriendshipAndReturnFriendInfo(friendRequestID uint64) *nexproto.FriendInfo {
	var senderPID uint32
	var recipientPID uint32

	if err := cassandraClusterSession.Query(`SELECT sender_pid, recipient_pid FROM pretendo_friends.friend_requests WHERE id=?`, friendRequestID).Scan(&senderPID, &recipientPID); err != nil {
		logger.Critical(err.Error())
		return nil
	}

	rand.Seed(time.Now().UnixNano())
	nodeID := rand.Intn(len(globals.SnowflakeNodes))

	snowflakeNode := globals.SnowflakeNodes[nodeID]

	friendshipID1 := uint64(snowflakeNode.Generate().Int64())
	friendshipID2 := uint64(snowflakeNode.Generate().Int64())

	acceptedTime := nex.NewDateTime(0)
	acceptedTime.FromTimestamp(time.Now())

	// Friendships are two-way relationships, not just one link between 2 entities
	// "A" has friend "B" and "B" has friend "A", so store both relationships

	if err := cassandraClusterSession.Query(`INSERT INTO pretendo_friends.friendships (id, user1_pid, user2_pid, date) VALUES (?, ?, ?, ?) IF NOT EXISTS`, friendshipID1, senderPID, recipientPID, acceptedTime.Value()).Exec(); err != nil {
		logger.Critical(err.Error())
		return nil
	}

	if err := cassandraClusterSession.Query(`INSERT INTO pretendo_friends.friendships (id, user1_pid, user2_pid, date) VALUES (?, ?, ?, ?) IF NOT EXISTS`, friendshipID2, recipientPID, senderPID, acceptedTime.Value()).Exec(); err != nil {
		logger.Critical(err.Error())
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
		if err := cassandraClusterSession.Query(`SELECT time FROM pretendo_friends.last_online WHERE pid=?`, senderPID).Scan(&lastOnlineTime); err != nil {
			if err == gocql.ErrNotFound {
				lastOnlineTime = nex.NewDateTime(0).Now()
			} else {
				logger.Critical(err.Error())
				lastOnlineTime = nex.NewDateTime(0).Now()
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
