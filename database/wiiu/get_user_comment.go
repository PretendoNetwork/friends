package database_wiiu

import (
	"database/sql"

	"github.com/PretendoNetwork/friends/database"
	"github.com/PretendoNetwork/nex-go"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/friends-wiiu/types"
)

// GetUserComment returns a user's comment
func GetUserComment(pid uint32) (*friends_wiiu_types.Comment, error) {
	comment := friends_wiiu_types.NewComment()
	comment.Unknown = 0

	var changed uint64 = 0

	err := database.Postgres.QueryRow(`SELECT comment, comment_changed FROM wiiu.user_data WHERE pid=$1`, pid).Scan(&comment.Contents, &changed)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, database.ErrPIDNotFound
		} else {
			return nil, err
		}
	}

	comment.LastChanged = nex.NewDateTime(changed)

	return comment, nil
}
