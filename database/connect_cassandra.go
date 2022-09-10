package database

import (
	"flag"
	"fmt"
	"time"

	"github.com/gocql/gocql"
)

var cluster *gocql.ClusterConfig
var cassandraClusterSession *gocql.Session

func connectCassandra() {
	// Connect to Cassandra

	var err error

	cluster = gocql.NewCluster("127.0.0.1")
	cluster.Timeout = 30 * time.Second

	createKeyspace("pretendo_friends")

	cluster.Keyspace = "pretendo_friends"

	cassandraClusterSession, err = cluster.CreateSession()

	if err != nil {
		logger.Critical(err.Error())
		return
	}

	// Create tables if missing

	if err := cassandraClusterSession.Query(`CREATE TABLE IF NOT EXISTS pretendo_friends.preferences (
		pid int PRIMARY KEY,
		show_online boolean,
		show_current_game boolean,
		block_friend_requests boolean
	)`).Exec(); err != nil {
		logger.Critical(err.Error())
		return
	}

	if err := cassandraClusterSession.Query(`CREATE TABLE IF NOT EXISTS pretendo_friends.blocks (
		id text PRIMARY KEY,
		blocker_pid int,
		blocked_pid int,
		date bigint
	)`).Exec(); err != nil {
		logger.Critical(err.Error())
		return
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
		logger.Critical(err.Error())
		return
	}

	if err := cassandraClusterSession.Query(`CREATE TABLE IF NOT EXISTS pretendo_friends.friendships (
		id bigint PRIMARY KEY,
		user1_pid int,
		user2_pid int,
		date bigint
	)`).Exec(); err != nil {
		logger.Critical(err.Error())
		return
	}

	if err := cassandraClusterSession.Query(`CREATE TABLE IF NOT EXISTS pretendo_friends.comments (
		pid int PRIMARY KEY,
		message text,
		changed bigint
	)`).Exec(); err != nil {
		logger.Critical(err.Error())
		return
	}

	if err := cassandraClusterSession.Query(`CREATE TABLE IF NOT EXISTS pretendo_friends.last_online (
		pid int PRIMARY KEY,
		time bigint
	)`).Exec(); err != nil {
		logger.Critical(err.Error())
		return
	}

	logger.Success("Connected to db")
}

// Adapted from gocql common_test.go
func createKeyspace(keyspace string) {
	flagRF := flag.Int("rf", 1, "replication factor for pretendo_friends keyspace")

	c := *cluster
	c.Keyspace = "system"
	c.Timeout = 30 * time.Second

	s, err := c.CreateSession()

	if err != nil {
		logger.Critical(err.Error())
	}

	defer s.Close()

	if err := s.Query(fmt.Sprintf(`CREATE KEYSPACE IF NOT EXISTS %s
	WITH replication = {
		'class' : 'SimpleStrategy',
		'replication_factor' : %d
	}`, keyspace, *flagRF)).Exec(); err != nil {
		logger.Critical(err.Error())
	}
}
