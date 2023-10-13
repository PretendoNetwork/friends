package nex_friends_wiiu

import (
	database_wiiu "github.com/PretendoNetwork/friends/database/wiiu"
	"github.com/PretendoNetwork/friends/globals"
	nex "github.com/PretendoNetwork/nex-go"
	friends_wiiu "github.com/PretendoNetwork/nex-protocols-go/friends-wiiu"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/friends-wiiu/types"
)

func UpdatePreference(err error, client *nex.Client, callID uint32, principalPreference *friends_wiiu_types.PrincipalPreference) uint32 {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nex.Errors.FPD.InvalidArgument
	}

	err = database_wiiu.UpdateUserPrincipalPreference(client.PID(), principalPreference)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return nex.Errors.FPD.Unknown
	}

	rmcResponse := nex.NewRMCResponse(friends_wiiu.ProtocolID, callID)
	rmcResponse.SetSuccess(friends_wiiu.MethodUpdatePreference, nil)

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewPacketV0(client, nil)

	responsePacket.SetVersion(0)
	responsePacket.SetSource(0xA1)
	responsePacket.SetDestination(0xAF)
	responsePacket.SetType(nex.DataPacket)
	responsePacket.SetPayload(rmcResponseBytes)

	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)

	globals.SecureServer.Send(responsePacket)

	return 0
}
