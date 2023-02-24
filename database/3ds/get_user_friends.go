package database_3ds

import (
	"database/sql"

	"github.com/PretendoNetwork/friends-secure/database"
	"github.com/PretendoNetwork/friends-secure/globals"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

// Get all of a user's friend relationships
func GetUserFriends(pid uint32) []*nexproto.FriendRelationship {
	friendRelationships := make([]*nexproto.FriendRelationship, 0)

	rows, err := database.Postgres.Query(`
	SELECT user2_pid, type FROM "3ds".friendships WHERE user1_pid=$1 AND type=1`, pid)
	if err != nil && err != sql.ErrNoRows {
		globals.Logger.Critical(err.Error())
	}

	for rows.Next() {
		relationship := nexproto.NewFriendRelationship()
		relationship.LFC = 0
		rows.Scan(&relationship.PID, &relationship.RelationshipType)

		friendRelationships = append(friendRelationships, relationship)
	}

	return friendRelationships
}
