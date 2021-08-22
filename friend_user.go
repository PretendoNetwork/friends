package main

import (
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

// FriendUser represents a user connected to the friends server
// Contains data relating only to friends
type FriendUser struct {
	NNID              string
	PID               uint32
	Comment           *nexproto.Comment
	FriendRequestsOut []*nexproto.FriendRequest
	FriendRequestsIn  []*nexproto.FriendRequest
	BlockedUsers      []*nexproto.BlacklistedPrincipal
	LastOnline        *nex.DateTime
	ActiveTitle       *nexproto.GameKey
	Notifications     []*nexproto.PersistentNotification
}

func (friendUser *FriendUser) FromPID(pid uint32) {

}

func NewFriendUser() *FriendUser {
	return &FriendUser{}
}
