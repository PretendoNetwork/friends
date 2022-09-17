package database_wiiu

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/PretendoNetwork/friends-secure/database"
	"github.com/PretendoNetwork/friends-secure/globals"
	"github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
	"github.com/gocql/gocql"
	"go.mongodb.org/mongo-driver/bson"
)

// Get a users friend list
func GetUserFriendList(pid uint32) []*nexproto.FriendInfo {
	friendList := make([]*nexproto.FriendInfo, 0)

	rows, err := database.Postgres.Query(`SELECT user2_pid, date FROM wiiu.friendships WHERE user1_pid=$1 AND active=true`, pid)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return friendList
	}

	for rows.Next() {
		var friendPID uint32
		var date uint64
		rows.Scan(&friendPID, &date)

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
			err := database.Postgres.QueryRow(`SELECT last_online FROM wiiu.user_data WHERE pid=$1`, friendPID).Scan(&lastOnlineTime)
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

		friendInfo.Status = GetUserComment(friendPID)
		friendInfo.BecameFriend = nex.NewDateTime(date)
		friendInfo.LastOnline = lastOnline
		friendInfo.Unknown = 0

		friendList = append(friendList, friendInfo)
	}

	return friendList
}
