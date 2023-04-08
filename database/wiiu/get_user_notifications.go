package database_wiiu

import friends_wiiu "github.com/PretendoNetwork/nex-protocols-go/friends/wiiu"

// Get notifications for a user
func GetUserNotifications(pid uint32) []*friends_wiiu.PersistentNotification {
	return make([]*friends_wiiu.PersistentNotification, 0)
}
