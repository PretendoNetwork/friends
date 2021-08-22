package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
	"github.com/gocql/gocql"
)

var cluster *gocql.ClusterConfig
var session *gocql.Session

func connectCassandra() {
	// Connect to Cassandra

	var err error

	cluster = gocql.NewCluster("127.0.0.1")
	cluster.Timeout = 30 * time.Second

	createKeyspace("pretendo_friends_wiiu")

	cluster.Keyspace = "pretendo_friends_wiiu"

	session, err = cluster.CreateSession()

	if err != nil {
		panic(err)
	}

	// Create tables if missing

	if err := session.Query(`CREATE TABLE IF NOT EXISTS pretendo_friends_wiiu.users (
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

	if err := session.Query(`CREATE TABLE IF NOT EXISTS pretendo_friends_wiiu.friendships (
		id bigint PRIMARY KEY,
		user1_pid int,
		user2_pid int,
		date bigint
	)`).Exec(); err != nil {
		log.Fatal(err)
	}

	if err := session.Query(`CREATE TABLE IF NOT EXISTS pretendo_friends_wiiu.blocks (
		id text PRIMARY KEY,
		blocker_pid int,
		blocked_pid int,
		date bigint
	)`).Exec(); err != nil {
		log.Fatal(err)
	}

	if err := session.Query(`CREATE TABLE IF NOT EXISTS pretendo_friends_wiiu.miis (
		pid int PRIMARY KEY,
		name text,
		unknown1 tinyint,
		unknown2 tinyint,
		data blob,
		date bigint
	)`).Exec(); err != nil {
		log.Fatal(err)
	}

	if err := session.Query(`CREATE TABLE IF NOT EXISTS pretendo_friends_wiiu.friend_requests (
		id bigint PRIMARY KEY,
		sender_pid int,
		recipient_pid int,
		sent_on bigint,
		message text
	)`).Exec(); err != nil {
		log.Fatal(err)
	}

	if err := session.Query(`CREATE TABLE IF NOT EXISTS pretendo_friends_wiiu.notifications (
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
	flagRF := flag.Int("rf", 1, "replication factor for pretendo_friends_wiiu keyspace")

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

// Update a users NNAInfo data
func updateNNAInfo(nnaInfo *nexproto.NNAInfo) {
	principalBasicInfo := nnaInfo.PrincipalBasicInfo

	userPID := principalBasicInfo.PID
	userNNID := principalBasicInfo.NNID
	userMii := principalBasicInfo.Mii

	// Insert users NNID into users table incase missing

	if err := session.Query(`UPDATE pretendo_friends_wiiu.users SET nnid = ? WHERE pid = ?`, userNNID, userPID).Exec(); err != nil {
		log.Fatal(err)
	}

	// Update user Mii data

	if err := session.Query(`UPDATE pretendo_friends_wiiu.miis SET
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

	if err := session.Query(`SELECT comment_message, comment_changed FROM pretendo_friends_wiiu.users WHERE pid = ? LIMIT 1`,
		pid).Consistency(gocql.One).Scan(&content, &changed); err != nil {
		return nil
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
