package notifications_wiiu

import (
	"github.com/PretendoNetwork/friends-secure/globals"
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func SendFriendshipRemoved(client *nex.Client, senderPID uint32) {
	nintendoNotificationEventGeneral := nexproto.NewNintendoNotificationEventGeneral()

	eventObject := nexproto.NewNintendoNotificationEvent()
	eventObject.Type = 26
	eventObject.SenderPID = senderPID
	eventObject.DataHolder = nex.NewDataHolder()
	eventObject.DataHolder.SetTypeName("NintendoNotificationEventGeneral")
	eventObject.DataHolder.SetObjectData(nintendoNotificationEventGeneral)

	stream := nex.NewStreamOut(globals.NEXServer)
	stream.WriteStructure(eventObject)

	rmcRequest := nex.NewRMCRequest()
	rmcRequest.SetProtocolID(nexproto.NintendoNotificationsProtocolID)
	rmcRequest.SetCallID(3810693103)
	rmcRequest.SetMethodID(nexproto.NintendoNotificationsMethodProcessNintendoNotificationEvent1)
	rmcRequest.SetParameters(stream.Bytes())

	rmcRequestBytes := rmcRequest.Bytes()

	requestPacket, _ := nex.NewPacketV0(client, nil)

	requestPacket.SetVersion(0)
	requestPacket.SetSource(0xA1)
	requestPacket.SetDestination(0xAF)
	requestPacket.SetType(nex.DataPacket)
	requestPacket.SetPayload(rmcRequestBytes)

	requestPacket.AddFlag(nex.FlagNeedsAck)
	requestPacket.AddFlag(nex.FlagReliable)

	globals.NEXServer.Send(requestPacket)
}
