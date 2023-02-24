package globals

import (
	"github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

type ConnectedUser struct {
	PID        uint32
	Platform   uint8
	Client     *nex.Client
	NNAInfo    *nexproto.NNAInfo
	Presence   *nexproto.NintendoPresence
	PresenceV2 *nexproto.NintendoPresenceV2
}

func NewConnectedUser() *ConnectedUser {
	return &ConnectedUser{}
}
