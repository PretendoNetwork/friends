package database_3ds

import (
	"database/sql"

	"github.com/PretendoNetwork/friends/database"
	"github.com/PretendoNetwork/nex-go/v2/types"
	friends_3ds_types "github.com/PretendoNetwork/nex-protocols-go/v2/friends-3ds/types"
)

// GetUserFriends returns all friend relationships of a user
func GetUserFriends(pid uint32) (*types.List[*friends_3ds_types.FriendRelationship], error) {
	friendRelationships := types.NewList[*friends_3ds_types.FriendRelationship]()
	friendRelationships.Type = friends_3ds_types.NewFriendRelationship()

	rows, err := database.Postgres.Query(`
	SELECT user2_pid, type FROM "3ds".friendships WHERE user1_pid=$1 AND type=1 LIMIT 100`, pid)
	if err != nil {
		if err == sql.ErrNoRows {
			return friendRelationships, database.ErrEmptyList
		} else {
			return friendRelationships, err
		}
	}

	for rows.Next() {
		var pid uint32
		var relationshipType uint8

		err := rows.Scan(&pid, &relationshipType)
		if err != nil {
			return friendRelationships, err
		}

		relationship := friends_3ds_types.NewFriendRelationship()

		relationship.LFC = types.NewPrimitiveU64(0)
		relationship.PID = types.NewPID(uint64(pid))
		relationship.RelationshipType = types.NewPrimitiveU8(relationshipType)

		friendRelationships.Append(relationship)
	}

	return friendRelationships, nil
}
