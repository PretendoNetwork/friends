package types

import (
	nex "github.com/PretendoNetwork/nex-go"
	friends_wiiu "github.com/PretendoNetwork/nex-protocols-go/friends-wiiu/types"
)

// FriendUser represents a user connected to the friends server
// Contains data relating only to friends
type FriendUser struct {
	NNID              string
	PID               uint32
	Comment           *friends_wiiu.Comment
	FriendRequestsOut []*friends_wiiu.FriendRequest
	FriendRequestsIn  []*friends_wiiu.FriendRequest
	BlockedUsers      []*friends_wiiu.BlacklistedPrincipal
	LastOnline        *nex.DateTime
	ActiveTitle       *friends_wiiu.GameKey
	Notifications     []*friends_wiiu.PersistentNotification
}

func (friendUser *FriendUser) FromPID(pid uint32) {

}

func NewFriendUser() *FriendUser {
	return &FriendUser{}
}
