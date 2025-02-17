package types

import (
	"github.com/PretendoNetwork/nex-go/v2"
	friends_3ds_types "github.com/PretendoNetwork/nex-protocols-go/v2/friends-3ds/types"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/v2/friends-wiiu/types"
)

type ConnectedUser struct {
	PID        uint32
	Platform   Platform
	Connection *nex.PRUDPConnection
	Presence   friends_3ds_types.NintendoPresence
	PresenceV2 friends_wiiu_types.NintendoPresenceV2
}

func NewConnectedUser() *ConnectedUser {
	return &ConnectedUser{
		Presence:   friends_3ds_types.NewNintendoPresence(),
		PresenceV2: friends_wiiu_types.NewNintendoPresenceV2(),
	}
}
