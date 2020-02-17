package main

import (
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func updateAndGetAllInformation(err error, client *nex.Client, callID uint32, nnaInfo *nexproto.NNAInfo, presence *nexproto.NintendoPresenceV2, birthday *nex.DateTime) {

	rmcResponseStream := nex.NewStreamOut(nexServer)

	// TODO: Make the following fields into structs and encode them

	comment := "Pretendo Servers"
	datetime := nex.NewDateTime(0)
	notificationString := "Test"

	//PrincipalPreference
	rmcResponseStream.WriteUInt8(0)
	rmcResponseStream.WriteUInt8(0)
	rmcResponseStream.WriteUInt8(0)

	//Comment
	rmcResponseStream.WriteUInt8(0)
	rmcResponseStream.WriteString(comment)
	rmcResponseStream.WriteUInt64LE(datetime.Now())
	//List<FriendInfo>
	rmcResponseStream.WriteUInt32LE(0)

	//List<FriendRequest> (Sent)
	rmcResponseStream.WriteUInt32LE(0)

	//List<FriendRequest> (Received)
	rmcResponseStream.WriteUInt32LE(0)

	//List<BlacklistedPrincipal>
	rmcResponseStream.WriteUInt32LE(0)

	//Unknown
	rmcResponseStream.WriteUInt8(0)

	//List<PersistentNotification>
	rmcResponseStream.WriteUInt32LE(1)

	//PersistentNotification
	rmcResponseStream.WriteUInt64LE(0xFFFFFFFFFFFFFFFF) //Unknown
	rmcResponseStream.WriteUInt32LE(0xFFFFFFFF)         //Unknown
	rmcResponseStream.WriteUInt32LE(0xFFFFFFFF)         //Unknown
	rmcResponseStream.WriteUInt32LE(0xFFFFFFFF)         //Unknown
	rmcResponseStream.WriteString(notificationString)   //Unknown

	//Unknown
	rmcResponseStream.WriteUInt8(0)

	rmcResponseBody := rmcResponseStream.Bytes()

	// Build response packet
	rmcResponse := nex.NewRMCResponse(nexproto.FriendsProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.FriendsMethodUpdateAndGetAllInformation, rmcResponseBody)

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
