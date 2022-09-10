package database

import "github.com/PretendoNetwork/nex-go"

func UpdateUserLastOnlineTime(pid uint32, lastOnline *nex.DateTime) {
	if err := cassandraClusterSession.Query(`UPDATE pretendo_friends.last_online SET time=? WHERE pid=?`, lastOnline.Value(), pid).Exec(); err != nil {
		logger.Critical(err.Error())
	}
}
