package types

import (
	"github.com/PretendoNetwork/nex-go"
	friends_3ds "github.com/PretendoNetwork/nex-protocols-go/friends/3ds"
	friends_wiiu "github.com/PretendoNetwork/nex-protocols-go/friends/wiiu"
)

type ConnectedUser struct {
	PID        uint32
	Platform   Platform
	Client     *nex.Client
	NNAInfo    *friends_wiiu.NNAInfo
	Presence   *friends_3ds.NintendoPresence
	PresenceV2 *friends_wiiu.NintendoPresenceV2
}

func NewConnectedUser() *ConnectedUser {
	return &ConnectedUser{}
}
