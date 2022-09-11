package database_wiiu

import nexproto "github.com/PretendoNetwork/nex-protocols-go"

// Get a users blacklist
func GetUserBlockList(pid uint32) []*nexproto.BlacklistedPrincipal {
	return make([]*nexproto.BlacklistedPrincipal, 0)
}
