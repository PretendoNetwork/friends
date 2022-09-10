package database

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/PretendoNetwork/friends-secure/globals"
	"github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
	"github.com/gocql/gocql"
	"go.mongodb.org/mongo-driver/bson"
)

// Get a users friend list
func GetUserFriendList(pid uint32) []*nexproto.FriendInfo {
	var sliceMap []map[string]interface{}
	var err error

	if sliceMap, err = cassandraClusterSession.Query(`SELECT user2_pid, date FROM pretendo_friends.friendships WHERE user1_pid=? ALLOW FILTERING`, pid).Iter().SliceMap(); err != nil {
		globals.Logger.Critical(err.Error())

		return make([]*nexproto.FriendInfo, 0)
	}

	friendList := make([]*nexproto.FriendInfo, 0)

	for i := 0; i < len(sliceMap); i++ {
		friendPID := uint32(sliceMap[i]["user2_pid"].(int))

		friendInfo := nexproto.NewFriendInfo()
		connectedUser := globals.ConnectedUsers[friendPID]
		var lastOnline *nex.DateTime

		if connectedUser != nil {
			// Online
			friendInfo.NNAInfo = connectedUser.NNAInfo
			friendInfo.Presence = connectedUser.Presence

			if friendInfo.NNAInfo == nil || friendInfo.NNAInfo.PrincipalBasicInfo == nil {
				// TODO: Fix this
				globals.Logger.Error(fmt.Sprintf("User %d has friend %d with bad presence data", pid, friendPID))
				if friendInfo.NNAInfo == nil {
					globals.Logger.Error(fmt.Sprintf("%d friendInfo.NNAInfo is nil", friendPID))
				} else {
					globals.Logger.Error(fmt.Sprintf("%d friendInfo.NNAInfo.PrincipalBasicInfo is nil", friendPID))
				}

				continue
			}

			lastOnline = nex.NewDateTime(0)
			lastOnline.FromTimestamp(time.Now())
		} else {
			// Offline
			friendUserInforation := GetUserInfoByPID(friendPID)
			encodedMiiData := friendUserInforation["mii"].(bson.M)["data"].(string)
			decodedMiiData, _ := base64.StdEncoding.DecodeString(encodedMiiData)

			friendInfo.NNAInfo = nexproto.NewNNAInfo()
			friendInfo.NNAInfo.PrincipalBasicInfo = nexproto.NewPrincipalBasicInfo()
			friendInfo.NNAInfo.PrincipalBasicInfo.PID = friendPID
			friendInfo.NNAInfo.PrincipalBasicInfo.NNID = friendUserInforation["username"].(string)
			friendInfo.NNAInfo.PrincipalBasicInfo.Mii = nexproto.NewMiiV2()
			friendInfo.NNAInfo.PrincipalBasicInfo.Mii.Name = friendUserInforation["mii"].(bson.M)["name"].(string)
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
			friendInfo.Presence.PID = 0
			friendInfo.Presence.GatheringID = 0
			friendInfo.Presence.ApplicationData = []byte{0x00}
			friendInfo.Presence.Unknown5 = 0
			friendInfo.Presence.Unknown6 = 0
			friendInfo.Presence.Unknown7 = 0

			var lastOnlineTime uint64
			if err := cassandraClusterSession.Query(`SELECT time FROM pretendo_friends.last_online WHERE pid=?`, friendPID).Scan(&lastOnlineTime); err != nil {
				if err == gocql.ErrNotFound {
					lastOnlineTime = nex.NewDateTime(0).Now()
				} else {
					globals.Logger.Critical(err.Error())
					lastOnlineTime = nex.NewDateTime(0).Now()
				}
			}

			lastOnline = nex.NewDateTime(lastOnlineTime) // TODO: Change this
		}

		friendInfo.Status = GetUserComment(friendPID)
		friendInfo.BecameFriend = nex.NewDateTime(uint64(sliceMap[i]["date"].(int64)))
		friendInfo.LastOnline = lastOnline
		friendInfo.Unknown = 0

		friendList = append(friendList, friendInfo)
	}

	return friendList
}
