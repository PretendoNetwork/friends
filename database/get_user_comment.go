package database

import (
	"github.com/PretendoNetwork/friends-secure/globals"
	"github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
	"github.com/gocql/gocql"
)

// Get a users comment
func GetUserComment(pid uint32) *nexproto.Comment {
	comment := nexproto.NewComment()
	comment.Unknown = 0

	var changed uint64 = 0

	if err := cassandraClusterSession.Query(`SELECT message, changed FROM pretendo_friends.comments WHERE pid=?`,
		pid).Consistency(gocql.One).Scan(&comment.Contents, &changed); err != nil {
		if err == gocql.ErrNotFound {
			comment.Contents = ""
		} else {
			comment.Contents = ""
			globals.Logger.Critical(err.Error())
		}
	}

	comment.LastChanged = nex.NewDateTime(changed)

	return comment
}
