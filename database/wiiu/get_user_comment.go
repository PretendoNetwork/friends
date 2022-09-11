package database_wiiu

import (
	"database/sql"

	"github.com/PretendoNetwork/friends-secure/database"
	"github.com/PretendoNetwork/friends-secure/globals"
	"github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

// Get a users comment
func GetUserComment(pid uint32) *nexproto.Comment {
	comment := nexproto.NewComment()
	comment.Unknown = 0

	var changed uint64 = 0

	err := database.Postgres.QueryRow(`SELECT comment, comment_changed FROM wiiu.user_data WHERE pid=$1`, pid).Scan(&comment.Contents, &changed)
	if err != nil {
		if err == sql.ErrNoRows {
			globals.Logger.Warning(err.Error())
		} else {
			globals.Logger.Critical(err.Error())
		}
	}

	comment.LastChanged = nex.NewDateTime(changed)

	return comment
}
