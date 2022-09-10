package globals

import (
	"github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

type ConnectedUser struct {
	PID      uint32
	Client   *nex.Client
	NNAInfo  *nexproto.NNAInfo
	Presence *nexproto.NintendoPresenceV2
}

func NewConnectedUser() *ConnectedUser {
	return &ConnectedUser{}
}
