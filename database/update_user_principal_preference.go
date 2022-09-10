package database

import (
	"github.com/PretendoNetwork/friends-secure/globals"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func UpdateUserPrincipalPreference(pid uint32, principalPreference *nexproto.PrincipalPreference) {
	if err := cassandraClusterSession.Query(`UPDATE pretendo_friends.preferences SET
		show_online=?,
		show_current_game=?,
		block_friend_requests=?
		WHERE pid=?`, principalPreference.ShowOnlinePresence, principalPreference.ShowCurrentTitle, principalPreference.BlockFriendRequests, pid).Exec(); err != nil {
		globals.Logger.Critical(err.Error())
	}
}
