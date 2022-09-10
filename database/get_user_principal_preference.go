package database

import nexproto "github.com/PretendoNetwork/nex-protocols-go"

func GetUserPrincipalPreference(pid uint32) *nexproto.PrincipalPreference {
	preference := nexproto.NewPrincipalPreference()

	_ = cassandraClusterSession.Query(`SELECT show_online, show_current_game, block_friend_requests FROM pretendo_friends.preferences WHERE pid=?`, pid).Scan(&preference.ShowOnlinePresence, &preference.ShowCurrentTitle, &preference.BlockFriendRequests)

	return preference
}
