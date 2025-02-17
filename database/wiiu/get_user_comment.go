package database_wiiu

import (
	"database/sql"

	"github.com/PretendoNetwork/friends/database"
	"github.com/PretendoNetwork/nex-go/v2/types"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/v2/friends-wiiu/types"
)

// GetUserComment returns a user's comment
func GetUserComment(pid uint32) (friends_wiiu_types.Comment, error) {
	comment := friends_wiiu_types.NewComment()

	var contents string
	var changed uint64 = 0

	row, err := database.Manager.QueryRow(`SELECT comment, comment_changed FROM wiiu.user_data WHERE pid=$1`, pid)
	if err != nil {
		return comment, err
	}

	err = row.Scan(&contents, &changed)
	if err != nil {
		if err == sql.ErrNoRows {
			return comment, database.ErrPIDNotFound
		} else {
			return comment, err
		}
	}

	comment.Contents = types.NewString(contents)
	comment.LastChanged = types.NewDateTime(changed)

	return comment, nil
}
