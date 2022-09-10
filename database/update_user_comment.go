package database

import (
	"github.com/PretendoNetwork/friends-secure/globals"
	"github.com/PretendoNetwork/nex-go"
)

// Update a users comment
func UpdateUserComment(pid uint32, message string) uint64 {
	changed := nex.NewDateTime(0).Now()

	if err := cassandraClusterSession.Query(`UPDATE pretendo_friends.comments SET message=?, changed=? WHERE pid=?`, message, changed, pid).Exec(); err != nil {
		globals.Logger.Critical(err.Error())
	}

	return changed
}
