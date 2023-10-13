package database_wiiu

import friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/friends-wiiu/types"

// GetUserNotifications returns notifications for a user
func GetUserNotifications(pid uint32) []*friends_wiiu_types.PersistentNotification {
	return make([]*friends_wiiu_types.PersistentNotification, 0)
}
