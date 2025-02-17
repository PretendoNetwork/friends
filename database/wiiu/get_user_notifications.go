package database_wiiu

import (
	"github.com/PretendoNetwork/nex-go/v2/types"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/v2/friends-wiiu/types"
)

// GetUserNotifications returns notifications for a user
func GetUserNotifications(pid uint32) types.List[friends_wiiu_types.PersistentNotification] {
	// TODO - Do this
	notifications := types.NewList[friends_wiiu_types.PersistentNotification]()

	return notifications
}
