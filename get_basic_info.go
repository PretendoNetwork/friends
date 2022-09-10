package main

import (
	"encoding/base64"

	"github.com/PretendoNetwork/friends-secure/database"
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
	"go.mongodb.org/mongo-driver/bson"
)

func getBasicInfo(err error, client *nex.Client, callID uint32, pids []uint32) {
	infos := make([]*nexproto.PrincipalBasicInfo, 0)

	for i := 0; i < len(pids); i++ {
		pid := pids[i]
		userInfo := database.GetUserInfoByPID(pid)

		info := nexproto.NewPrincipalBasicInfo()
		info.PID = pid
		info.NNID = userInfo["username"].(string)
		info.Mii = nexproto.NewMiiV2()
		info.Unknown = 2 // idk

		encodedMiiData := userInfo["mii"].(bson.M)["data"].(string)
		decodedMiiData, _ := base64.StdEncoding.DecodeString(encodedMiiData)

		info.Mii.Name = userInfo["mii"].(bson.M)["name"].(string)
		info.Mii.Unknown1 = 0
		info.Mii.Unknown2 = 0
		info.Mii.Data = decodedMiiData
		info.Mii.Datetime = nex.NewDateTime(0)

		infos = append(infos, info)
	}

	rmcResponseStream := nex.NewStreamOut(nexServer)

	rmcResponseStream.WriteListStructure(infos)

	rmcResponseBody := rmcResponseStream.Bytes()

	// Build response packet
	rmcResponse := nex.NewRMCResponse(nexproto.FriendsWiiUProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.FriendsWiiUMethodGetBasicInfo, rmcResponseBody)

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewPacketV0(client, nil)

	responsePacket.SetVersion(0)
	responsePacket.SetSource(0xA1)
	responsePacket.SetDestination(0xAF)
	responsePacket.SetType(nex.DataPacket)
	responsePacket.SetPayload(rmcResponseBytes)

	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)

	nexServer.Send(responsePacket)
}
