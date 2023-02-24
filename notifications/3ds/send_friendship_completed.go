package notifications_3ds

import (
	"github.com/PretendoNetwork/friends-secure/globals"
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func SendFriendshipCompleted(client *nex.Client, friendPID uint32, senderPID uint32) {
	notificationEvent := nexproto.NewNintendoNotificationEventGeneral()
	notificationEvent.U32Param = 0
	notificationEvent.U64Param1 = 0
	notificationEvent.U64Param2 = uint64(friendPID)

	eventObject := nexproto.NewNintendoNotificationEvent()
	eventObject.Type = 7
	eventObject.SenderPID = senderPID
	eventObject.DataHolder = nex.NewDataHolder()
	eventObject.DataHolder.SetTypeName("NintendoNotificationEventGeneral")
	eventObject.DataHolder.SetObjectData(notificationEvent)

	stream := nex.NewStreamOut(globals.NEXServer)
	eventObjectBytes := eventObject.Bytes(stream)

	rmcRequest := nex.NewRMCRequest()
	rmcRequest.SetProtocolID(nexproto.NintendoNotificationsProtocolID)
	rmcRequest.SetCallID(3810693103)
	rmcRequest.SetMethodID(nexproto.NintendoNotificationsMethodProcessNintendoNotificationEvent1)
	rmcRequest.SetParameters(eventObjectBytes)

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
