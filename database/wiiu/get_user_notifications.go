package database_wiiu

import nexproto "github.com/PretendoNetwork/nex-protocols-go"

// Get notifications for a user
func GetUserNotifications(pid uint32) []*nexproto.PersistentNotification {
	return make([]*nexproto.PersistentNotification, 0)
}
