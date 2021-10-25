package main

import (
	"context"
	"flag"
	"fmt"
	"log"
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

	if err := cassandraClusterSession.Query(`CREATE TABLE IF NOT EXISTS pretendo_friends.preferences (
		pid int PRIMARY KEY,
		show_online boolean,
		show_current_game boolean,
		block_friend_requests boolean
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

	if err := cassandraClusterSession.Query(`CREATE TABLE IF NOT EXISTS pretendo_friends.blocks (
		id text PRIMARY KEY,
		blocker_pid int,
		blocked_pid int,
		date bigint
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

	if err := cassandraClusterSession.Query(`CREATE TABLE IF NOT EXISTS pretendo_friends.friend_requests (
		id bigint PRIMARY KEY,
		sender_pid int,
		recipient_pid int,
		sent_on bigint,
		message text
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
	return make([]*nexproto.FriendInfo, 0)
}

// Get a users sent friend requests
func getUserFriendRequestsOut(pid uint32) {}

// Get a users received friend requests
func getUserFriendRequestsIn(pid uint32) {}

// Get a users blacklist
func getUserBlockList(pid uint32) {}

// Get notifications for a user
func getUserNotifications(pid uint32) {}

func updateUserPrincipalPreference(pid uint32, principalPreference *nexproto.PrincipalPreference) {
	if err := cassandraClusterSession.Query(`UPDATE pretendo_friends.preferences SET show_online=?, show_current_game=?, block_friend_requests=? WHERE pid=?`, principalPreference.Unknown1, principalPreference.Unknown2, principalPreference.Unknown3, pid).Exec(); err != nil {
		log.Fatal(err)
	}
}

func getUserPrincipalPreference(pid uint32) *nexproto.PrincipalPreference {
	var showOnline bool
	var showCurrentGame bool
	var blockFriendRequests bool

	_ = cassandraClusterSession.Query(`SELECT show_online, show_current_game, block_friend_requests FROM pretendo_friends.preferences WHERE pid=?`, pid).Scan(&showOnline, &showCurrentGame, &blockFriendRequests)

	preference := nexproto.NewPrincipalPreference()
	preference.Unknown1 = showOnline
	preference.Unknown2 = showCurrentGame
	preference.Unknown3 = blockFriendRequests

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
