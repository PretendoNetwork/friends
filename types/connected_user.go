package types

import (
	"github.com/PretendoNetwork/nex-go"
	friends_3ds_types "github.com/PretendoNetwork/nex-protocols-go/friends-3ds/types"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/friends-wiiu/types"
)

type ConnectedUser struct {
	PID        uint32
	Platform   Platform
	Client     *nex.Client
	NNAInfo    *friends_wiiu_types.NNAInfo
	Presence   *friends_3ds_types.NintendoPresence
	PresenceV2 *friends_wiiu_types.NintendoPresenceV2
}

func NewConnectedUser() *ConnectedUser {
	return &ConnectedUser{}
}
