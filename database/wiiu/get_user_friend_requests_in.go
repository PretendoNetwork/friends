package database_wiiu

import (
	"encoding/base64"

	"github.com/PretendoNetwork/friends-secure/database"
	"github.com/PretendoNetwork/friends-secure/globals"
	"github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
	"go.mongodb.org/mongo-driver/bson"
)

// Get a users received friend requests
func GetUserFriendRequestsIn(pid uint32) []*nexproto.FriendRequest {
	friendRequestsIn := make([]*nexproto.FriendRequest, 0)

	rows, err := database.Postgres.Query(`SELECT id, sender_pid, sent_on, expires_on, message, received FROM wiiu.friend_requests WHERE recipient_pid=$1 AND accepted=false AND denied=false`, pid)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return friendRequestsIn
	}

	for rows.Next() {
		var id uint64
		var senderPID uint32
		var sentOn uint64
		var expiresOn uint64
		var message string
		var received bool
		rows.Scan(&id, &senderPID, &sentOn, &expiresOn, &message, &received)

		senderUserInforation := GetUserInfoByPID(senderPID)
		encodedMiiData := senderUserInforation["mii"].(bson.M)["data"].(string)
		decodedMiiData, _ := base64.StdEncoding.DecodeString(encodedMiiData)

		friendRequest := nexproto.NewFriendRequest()

		friendRequest.PrincipalInfo = nexproto.NewPrincipalBasicInfo()
		friendRequest.PrincipalInfo.PID = senderPID
		friendRequest.PrincipalInfo.NNID = senderUserInforation["username"].(string)
		friendRequest.PrincipalInfo.Mii = nexproto.NewMiiV2()
		friendRequest.PrincipalInfo.Mii.Name = senderUserInforation["mii"].(bson.M)["name"].(string)
		friendRequest.PrincipalInfo.Mii.Unknown1 = 0 // replaying from real server
		friendRequest.PrincipalInfo.Mii.Unknown2 = 0 // replaying from real server
		friendRequest.PrincipalInfo.Mii.Data = decodedMiiData
		friendRequest.PrincipalInfo.Mii.Datetime = nex.NewDateTime(0)
		friendRequest.PrincipalInfo.Unknown = 2 // replaying from real server

		friendRequest.Message = nexproto.NewFriendRequestMessage()
		friendRequest.Message.FriendRequestID = id
		friendRequest.Message.Received = received
		friendRequest.Message.Unknown2 = 1
		friendRequest.Message.Message = message
		friendRequest.Message.Unknown3 = 0
		friendRequest.Message.Unknown4 = ""
		friendRequest.Message.GameKey = nexproto.NewGameKey()
		friendRequest.Message.GameKey.TitleID = 0
		friendRequest.Message.GameKey.TitleVersion = 0
		friendRequest.Message.Unknown5 = nex.NewDateTime(134222053376) // idk what this value means but its always this
		friendRequest.Message.ExpiresOn = nex.NewDateTime(expiresOn)
		friendRequest.SentOn = nex.NewDateTime(sentOn)

		friendRequestsIn = append(friendRequestsIn, friendRequest)
	}

	return friendRequestsIn
}
