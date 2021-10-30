package main

import (
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
	"github.com/gocql/gocql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var cluster *gocql.ClusterConfig
var cassandraClusterSession *gocql.Session

var mongoClient *mongo.Client
var mongoContext context.Context
var mongoDatabase *mongo.Database
var mongoCollection *mongo.Collection

func connectMongo() {
	mongoClient, _ = mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	mongoContext, _ = context.WithTimeout(context.Background(), 10*time.Second)
	_ = mongoClient.Connect(mongoContext)

	mongoDatabase = mongoClient.Database("pretendo")
	mongoCollection = mongoDatabase.Collection("pnids")
}

func connectCassandra() {
	// Connect to Cassandra

	var err error

	cluster = gocql.NewCluster("127.0.0.1")
	cluster.Timeout = 30 * time.Second

	createKeyspace("pretendo_friends")

	cluster.Keyspace = "pretendo_friends"

	cassandraClusterSession, err = cluster.CreateSession()

	if err != nil {
		panic(err)
	}

	// Create tables if missing

	if err := cassandraClusterSession.Query(`CREATE TABLE IF NOT EXISTS pretendo_friends.preferences (
		pid int PRIMARY KEY,
		show_online boolean,
		show_current_game boolean,
		block_friend_requests boolean
	)`).Exec(); err != nil {
		log.Fatal(err)
	}

	if err := cassandraClusterSession.Query(`CREATE TABLE IF NOT EXISTS pretendo_friends.blocks (
		id text PRIMARY KEY,
		blocker_pid int,
		blocked_pid int,
		date bigint
	)`).Exec(); err != nil {
		log.Fatal(err)
	}

	if err := cassandraClusterSession.Query(`CREATE TABLE IF NOT EXISTS pretendo_friends.friend_requests (
		id bigint PRIMARY KEY,
		sender_pid int,
		recipient_pid int,
		sent_on bigint,
		expires_on bigint,
		message text,
		received boolean,
		accepted boolean,
		denied boolean
	)`).Exec(); err != nil {
		log.Fatal(err)
	}

	if err := cassandraClusterSession.Query(`CREATE TABLE IF NOT EXISTS pretendo_friends.friendships (
		id bigint PRIMARY KEY,
		user1_pid int,
		user2_pid int,
		date bigint
	)`).Exec(); err != nil {
		log.Fatal(err)
	}

	if err := cassandraClusterSession.Query(`CREATE TABLE IF NOT EXISTS pretendo_friends.users (
		pid int PRIMARY KEY,
		nnid text,
		changed_flags int,
		comment_message text,
		comment_changed bigint,
		last_online bigint,
		gathering_id int,
		active_title_id bigint,
		active_title_version smallint,
		active_title_game_server_id int,
		active_title_data blob,
		friendships list<int>,
		blocks list<int>
	)`).Exec(); err != nil {
		log.Fatal(err)
	}

	if err := cassandraClusterSession.Query(`CREATE TABLE IF NOT EXISTS pretendo_friends.miis (
		pid int PRIMARY KEY,
		name text,
		unknown1 tinyint,
		unknown2 tinyint,
		data blob,
		date bigint
	)`).Exec(); err != nil {
		log.Fatal(err)
	}

	if err := cassandraClusterSession.Query(`CREATE TABLE IF NOT EXISTS pretendo_friends.notifications (
		id bigint PRIMARY KEY,
		sender_pid int,
		recipient_pid int,
		event_type int,
		message text
	)`).Exec(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to Cassandra")
}

// Adapted from gocql common_test.go
func createKeyspace(keyspace string) {
	flagRF := flag.Int("rf", 1, "replication factor for pretendo_friends keyspace")

	c := *cluster
	c.Keyspace = "system"
	c.Timeout = 30 * time.Second

	s, err := c.CreateSession()

	if err != nil {
		panic(err)
	}

	defer s.Close()

	if err := s.Query(fmt.Sprintf(`CREATE KEYSPACE IF NOT EXISTS %s
	WITH replication = {
		'class' : 'SimpleStrategy',
		'replication_factor' : %d
	}`, keyspace, *flagRF)).Exec(); err != nil {
		log.Fatal(err)
	}
}

////////////////////////////////
//                            //
// Cassandra database methods //
//                            //
////////////////////////////////

// Update a users NNAInfo data
func updateNNAInfo(nnaInfo *nexproto.NNAInfo) {
	principalBasicInfo := nnaInfo.PrincipalBasicInfo

	userPID := principalBasicInfo.PID
	userNNID := principalBasicInfo.NNID
	userMii := principalBasicInfo.Mii

	// Insert users NNID into users table incase missing

	if err := cassandraClusterSession.Query(`UPDATE pretendo_friends.users SET nnid = ? WHERE pid = ?`, userNNID, userPID).Exec(); err != nil {
		log.Fatal(err)
	}

	// Update user Mii data

	if err := cassandraClusterSession.Query(`UPDATE pretendo_friends.miis SET
		data = ?,
		name = ?,
		date = ?
		WHERE pid = ?`,
		userMii.Data, userMii.Name, userMii.Datetime.Value(), userPID).Exec(); err != nil {
		log.Fatal(err)
	}
}

// Update a users NintendoPresenceV2 data
func updateNintendoPresenceV2(presence *nexproto.NintendoPresenceV2) {

}

// Get a users comment
func getUserComment(pid uint32) *nexproto.Comment {
	var content string
	var changed uint64

	if err := cassandraClusterSession.Query(`SELECT comment_message, comment_changed FROM pretendo_friends.users WHERE pid = ? LIMIT 1`,
		pid).Consistency(gocql.One).Scan(&content, &changed); err != nil {
		comment := nexproto.NewComment()
		comment.Unknown = 0
		comment.Contents = ""
		comment.LastChanged = nex.NewDateTime(0)

		return comment
		// TODO: Handle the error
	}

	comment := nexproto.NewComment()

	comment.Unknown = 0
	comment.Contents = content
	comment.LastChanged = nex.NewDateTime(changed)

	return comment
}

// Get a users friend list
func getUserFriendList(pid uint32) []*nexproto.FriendInfo {
	var sliceMap []map[string]interface{}
	var err error

	if sliceMap, err = cassandraClusterSession.Query(`SELECT user2_pid, date FROM pretendo_friends.friendships WHERE user1_pid=? ALLOW FILTERING`, pid).Iter().SliceMap(); err != nil {
		log.Fatal(err)
	}

	friendList := make([]*nexproto.FriendInfo, 0)

	for i := 0; i < len(sliceMap); i++ {
		friendPID := uint32(sliceMap[i]["user2_pid"].(int))

		friendInfo := nexproto.NewFriendInfo()
		connectedUser := connectedUsers[friendPID]

		if connectedUser != nil {
			// Online
			friendInfo.NNAInfo = connectedUser.NNAInfo
			friendInfo.Presence = connectedUser.Presence
		} else {
			// Offline
			friendUserInforation := getUserInfoByPID(friendPID)
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
		}

		friendInfo.Status = getUserComment(friendPID)
		friendInfo.BecameFriend = nex.NewDateTime(uint64(sliceMap[i]["date"].(int64)))
		friendInfo.LastOnline = nex.NewDateTime(uint64(sliceMap[i]["date"].(int64))) // TODO: Change this
		friendInfo.Unknown = 0

		friendList = append(friendList, friendInfo)
	}

	return friendList
}

// Get a users sent friend requests
func getUserFriendRequestsOut(pid uint32) []*nexproto.FriendRequest {
	var sliceMap []map[string]interface{}
	var err error

	if sliceMap, err = cassandraClusterSession.Query(`SELECT id, recipient_pid, sent_on, expires_on, message, received FROM pretendo_friends.friend_requests WHERE sender_pid=? AND accepted=false AND denied=false ALLOW FILTERING`, pid).Iter().SliceMap(); err != nil {
		log.Fatal(err)
	}

	friendRequestsOut := make([]*nexproto.FriendRequest, 0)

	for i := 0; i < len(sliceMap); i++ {
		recipientPID := uint32(sliceMap[i]["recipient_pid"].(int))

		recipientUserInforation := getUserInfoByPID(recipientPID)
		encodedMiiData := recipientUserInforation["mii"].(bson.M)["data"].(string)
		decodedMiiData, _ := base64.StdEncoding.DecodeString(encodedMiiData)

		friendRequest := nexproto.NewFriendRequest()

		friendRequest.PrincipalInfo = nexproto.NewPrincipalBasicInfo()
		friendRequest.PrincipalInfo.PID = recipientPID
		friendRequest.PrincipalInfo.NNID = recipientUserInforation["username"].(string)
		friendRequest.PrincipalInfo.Mii = nexproto.NewMiiV2()
		friendRequest.PrincipalInfo.Mii.Name = recipientUserInforation["mii"].(bson.M)["name"].(string)
		friendRequest.PrincipalInfo.Mii.Unknown1 = 0 // replaying from real server
		friendRequest.PrincipalInfo.Mii.Unknown2 = 0 // replaying from real server
		friendRequest.PrincipalInfo.Mii.Data = decodedMiiData
		friendRequest.PrincipalInfo.Mii.Datetime = nex.NewDateTime(0)
		friendRequest.PrincipalInfo.Unknown = 2 // replaying from real server

		friendRequest.Message = nexproto.NewFriendRequestMessage()
		friendRequest.Message.FriendRequestID = uint64(sliceMap[i]["id"].(int64))
		friendRequest.Message.Received = sliceMap[i]["received"].(bool)
		friendRequest.Message.Unknown2 = 1
		friendRequest.Message.Message = sliceMap[i]["message"].(string)
		friendRequest.Message.Unknown3 = 0
		friendRequest.Message.Unknown4 = ""
		friendRequest.Message.GameKey = nexproto.NewGameKey()
		friendRequest.Message.GameKey.TitleID = 0
		friendRequest.Message.GameKey.TitleVersion = 0
		friendRequest.Message.Unknown5 = nex.NewDateTime(134222053376) // idk what this value means but its always this
		friendRequest.Message.ExpiresOn = nex.NewDateTime(uint64(sliceMap[i]["expires_on"].(int64)))
		friendRequest.SentOn = nex.NewDateTime(uint64(sliceMap[i]["sent_on"].(int64)))

		friendRequestsOut = append(friendRequestsOut, friendRequest)
	}

	return friendRequestsOut
}

// Get a users received friend requests
func getUserFriendRequestsIn(pid uint32) []*nexproto.FriendRequest {
	var sliceMap []map[string]interface{}
	var err error

	if sliceMap, err = cassandraClusterSession.Query(`SELECT id, sender_pid, sent_on, expires_on, message, received FROM pretendo_friends.friend_requests WHERE recipient_pid=? AND accepted=false AND denied=false ALLOW FILTERING`, pid).Iter().SliceMap(); err != nil {
		log.Fatal(err)
	}

	friendRequestsOut := make([]*nexproto.FriendRequest, 0)

	for i := 0; i < len(sliceMap); i++ {
		senderPID := uint32(sliceMap[i]["sender_pid"].(int))

		senderUserInforation := getUserInfoByPID(senderPID)
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
		friendRequest.Message.FriendRequestID = uint64(sliceMap[i]["id"].(int64))
		friendRequest.Message.Received = sliceMap[i]["received"].(bool)
		friendRequest.Message.Unknown2 = 1
		friendRequest.Message.Message = sliceMap[i]["message"].(string)
		friendRequest.Message.Unknown3 = 0
		friendRequest.Message.Unknown4 = ""
		friendRequest.Message.GameKey = nexproto.NewGameKey()
		friendRequest.Message.GameKey.TitleID = 0
		friendRequest.Message.GameKey.TitleVersion = 0
		friendRequest.Message.Unknown5 = nex.NewDateTime(134222053376) // idk what this value means but its always this
		friendRequest.Message.ExpiresOn = nex.NewDateTime(uint64(sliceMap[i]["expires_on"].(int64)))
		friendRequest.SentOn = nex.NewDateTime(uint64(sliceMap[i]["sent_on"].(int64)))

		friendRequestsOut = append(friendRequestsOut, friendRequest)
	}

	return friendRequestsOut
}

// Get a users blacklist
func getUserBlockList(pid uint32) {}

// Get notifications for a user
func getUserNotifications(pid uint32) {}

func updateUserPrincipalPreference(pid uint32, principalPreference *nexproto.PrincipalPreference) {
	if err := cassandraClusterSession.Query(`UPDATE pretendo_friends.preferences SET
		show_online=?,
		show_current_game=?,
		block_friend_requests=?
		WHERE pid=?`, principalPreference.ShowOnlinePresence, principalPreference.ShowCurrentTitle, principalPreference.BlockFriendRequests, pid).Exec(); err != nil {
		log.Fatal(err)
	}
}

func getUserPrincipalPreference(pid uint32) *nexproto.PrincipalPreference {
	preference := nexproto.NewPrincipalPreference()

	_ = cassandraClusterSession.Query(`SELECT show_online, show_current_game, block_friend_requests FROM pretendo_friends.preferences WHERE pid=?`, pid).Scan(&preference.ShowOnlinePresence, &preference.ShowCurrentTitle, &preference.BlockFriendRequests)

	return preference
}

func isFriendRequestBlocked(requesterPID uint32, requestedPID uint32) bool {
	if err := cassandraClusterSession.Query(`SELECT id FROM pretendo_friends.blocks WHERE blocker_pid=? AND blocked_pid=? LIMIT 1 ALLOW FILTERING`, requestedPID, requesterPID).Scan(); err != nil {
		if err == gocql.ErrNotFound {
			// Assume no block record was found
			return false
		}

		// TODO: Error handling
	}

	// Assume a block record was found
	return true
}

func saveFriendRequest(friendRequestID uint64, senderPID uint32, recipientPID uint32, sentTime uint64, expireTime uint64, message string) {
	if err := cassandraClusterSession.Query(`INSERT INTO pretendo_friends.friend_requests (id, sender_pid, recipient_pid, sent_on, expires_on, message, received, accepted, denied) VALUES (?, ?, ?, ?, ?, ?, false, false, false) IF NOT EXISTS`, friendRequestID, senderPID, recipientPID, sentTime, expireTime, message).Exec(); err != nil {
		log.Fatal(err)
	}
}

func setFriendRequestReceived(friendRequestID uint64) {
	if err := cassandraClusterSession.Query(`UPDATE pretendo_friends.friend_requests SET received=true WHERE id=?`, friendRequestID).Exec(); err != nil {
		log.Fatal(err)
	}
}

func setFriendRequestAccepted(friendRequestID uint64) {
	if err := cassandraClusterSession.Query(`UPDATE pretendo_friends.friend_requests SET accepted=true WHERE id=?`, friendRequestID).Exec(); err != nil {
		log.Fatal(err)
	}
}

func acceptFriendshipAndReturnFriendInfo(friendRequestID uint64) *nexproto.FriendInfo {
	var senderPID uint32
	var recipientPID uint32

	if err := cassandraClusterSession.Query(`SELECT sender_pid, recipient_pid FROM pretendo_friends.friend_requests WHERE id=?`, friendRequestID).Scan(&senderPID, &recipientPID); err != nil {
		log.Fatal(err)
	}

	rand.Seed(time.Now().UnixNano())
	nodeID := rand.Intn(len(snowflakeNodes))

	snowflakeNode := snowflakeNodes[nodeID]

	friendshipID1 := uint64(snowflakeNode.Generate().Int64())
	friendshipID2 := uint64(snowflakeNode.Generate().Int64())

	acceptedTime := nex.NewDateTime(0)
	acceptedTime.FromTimestamp(time.Now())

	// Friendships are two-way relationships, not just one link between 2 entities
	// "A" has friend "B" and "B" has friend "A", so store both relationships

	if err := cassandraClusterSession.Query(`INSERT INTO pretendo_friends.friendships (id, user1_pid, user2_pid, date) VALUES (?, ?, ?, ?) IF NOT EXISTS`, friendshipID1, senderPID, recipientPID, acceptedTime.Value()).Exec(); err != nil {
		log.Fatal(err)
	}

	if err := cassandraClusterSession.Query(`INSERT INTO pretendo_friends.friendships (id, user1_pid, user2_pid, date) VALUES (?, ?, ?, ?) IF NOT EXISTS`, friendshipID2, recipientPID, senderPID, acceptedTime.Value()).Exec(); err != nil {
		log.Fatal(err)
	}

	setFriendRequestAccepted(friendRequestID)

	friendInfo := nexproto.NewFriendInfo()
	connectedUser := connectedUsers[senderPID]

	if connectedUser != nil {
		// Online
		friendInfo.NNAInfo = connectedUser.NNAInfo
		friendInfo.Presence = connectedUser.Presence
	} else {
		// Offline
		senderUserInforation := getUserInfoByPID(senderPID)
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
		friendInfo.Presence.PID = 0
		friendInfo.Presence.GatheringID = 0
		friendInfo.Presence.ApplicationData = []byte{0x00}
		friendInfo.Presence.Unknown5 = 0
		friendInfo.Presence.Unknown6 = 0
		friendInfo.Presence.Unknown7 = 0
	}

	friendInfo.Status = getUserComment(senderPID)
	friendInfo.BecameFriend = acceptedTime
	friendInfo.LastOnline = acceptedTime // TODO: Change this
	friendInfo.Unknown = 0

	return friendInfo
}

//////////////////////////////
//                          //
// MongoDB database methods //
//                          //
//////////////////////////////

func getUserInfoByPID(pid uint32) bson.M {
	var result bson.M

	err := mongoCollection.FindOne(context.TODO(), bson.D{{Key: "pid", Value: pid}}, options.FindOne()).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}

		panic(err)
	}

	return result
}
