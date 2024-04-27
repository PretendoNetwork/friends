package database_wiiu

import (
	"github.com/PretendoNetwork/friends/database"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/v2/friends-wiiu/types"
)

// UpdateUserPrincipalPreference updates the user preferences
func UpdateUserPrincipalPreference(pid uint32, principalPreference *friends_wiiu_types.PrincipalPreference) error {
	_, err := database.Manager.Exec(`
		INSERT INTO wiiu.user_data (pid, show_online, show_current_game, block_friend_requests)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (pid)
		DO UPDATE SET 
		show_online = $2,
		show_current_game = $3,
		block_friend_requests = $4`, pid, principalPreference.ShowOnlinePresence.Value, principalPreference.ShowCurrentTitle.Value, principalPreference.BlockFriendRequests.Value)

	if err != nil {
		return err
	}

	return nil
}
