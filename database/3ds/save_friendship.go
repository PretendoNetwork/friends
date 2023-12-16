package database_3ds

import (
	"github.com/PretendoNetwork/friends/database"
	"github.com/PretendoNetwork/nex-go"
	friends_3ds_types "github.com/PretendoNetwork/nex-protocols-go/friends-3ds/types"
)

// SaveFriendship saves a friend relationship for a user
func SaveFriendship(senderPID uint32, recipientPID uint32) (*friends_3ds_types.FriendRelationship, error) {
	friendRelationship := friends_3ds_types.NewFriendRelationship()
	friendRelationship.PID = nex.NewPID(recipientPID)
	friendRelationship.LFC = 0
	friendRelationship.RelationshipType = 0 // * Incomplete

	// * Ensure that we inputted a valid user.
	var found bool
	err := database.Postgres.QueryRow(`SELECT COUNT(*) FROM "3ds".user_data WHERE pid=$1 LIMIT 1`, recipientPID).Scan(&found)
	if err != nil {
		return nil, err
	}
	if !found {
		friendRelationship.RelationshipType = 2 // * Non-existent
		return friendRelationship, nil
	}

	// * Get the other side's relationship, we need to know if we've already got one sent to us.
	err = database.Postgres.QueryRow(`
	SELECT COUNT(*) FROM "3ds".friendships WHERE user1_pid=$1 AND user2_pid=$2 AND type=0 LIMIT 1`, recipientPID, senderPID).Scan(&found)
	if err != nil {
		return nil, err
	}
	if !found {
		_, err = database.Postgres.Exec(`
		INSERT INTO "3ds".friendships (user1_pid, user2_pid, type)
		VALUES ($1, $2, 0)
		ON CONFLICT (user1_pid, user2_pid)
		DO NOTHING`, senderPID, recipientPID)
		if err != nil {
			return nil, err
		}
		return friendRelationship, nil
	}

	acceptedTime := nex.NewDateTime(0).Now().Value()

	// * We need to have two relationships for both sides as friend relationships are not one single object.
	_, err = database.Postgres.Exec(`
		INSERT INTO "3ds".friendships (user1_pid, user2_pid, date, type)
		VALUES ($1, $2, $3, 1)
		ON CONFLICT (user1_pid, user2_pid)
		DO UPDATE SET
		date = $3,
		type = 1`, senderPID, recipientPID, acceptedTime)
	if err != nil {
		return nil, err
	}

	_, err = database.Postgres.Exec(`
		INSERT INTO "3ds".friendships (user1_pid, user2_pid, date, type)
		VALUES ($1, $2, $3, 1)
		ON CONFLICT (user1_pid, user2_pid)
		DO UPDATE SET
		date = $3,
		type = 1`, recipientPID, senderPID, acceptedTime)
	if err != nil {
		return nil, err
	}

	friendRelationship.RelationshipType = 1 // * Complete
	return friendRelationship, nil
}
