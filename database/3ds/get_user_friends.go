package database_3ds

import (
	"database/sql"

	"github.com/PretendoNetwork/friends/database"
	friends_3ds_types "github.com/PretendoNetwork/nex-protocols-go/friends-3ds/types"
)

// GetUserFriends returns all friend relationships of a user
func GetUserFriends(pid uint32) ([]*friends_3ds_types.FriendRelationship, error) {
	friendRelationships := make([]*friends_3ds_types.FriendRelationship, 0)

	rows, err := database.Postgres.Query(`
	SELECT user2_pid, type FROM "3ds".friendships WHERE user1_pid=$1 AND type=1 LIMIT 100`, pid)
	if err != nil {
		if err ==  sql.ErrNoRows {
			return friendRelationships, database.ErrEmptyList
		} else {
			return friendRelationships, err
		}
	}

	for rows.Next() {
		relationship := friends_3ds_types.NewFriendRelationship()
		relationship.LFC = 0
		rows.Scan(&relationship.PID, &relationship.RelationshipType)

		friendRelationships = append(friendRelationships, relationship)
	}

	return friendRelationships, nil
}
