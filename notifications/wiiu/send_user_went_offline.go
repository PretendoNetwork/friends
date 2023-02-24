package notifications_wiiu

import (
	"time"

	database_wiiu "github.com/PretendoNetwork/friends-secure/database/wiiu"
	"github.com/PretendoNetwork/friends-secure/globals"
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func SendUserWentOfflineGlobally(client *nex.Client) {
	friendsList := database_wiiu.GetUserFriendList(client.PID())

	for i := 0; i < len(friendsList); i++ {
		SendUserWentOffline(client, friendsList[i].NNAInfo.PrincipalBasicInfo.PID)
	}
}

func SendUserWentOffline(client *nex.Client, pid uint32) {
	lastOnline := nex.NewDateTime(0)
	lastOnline.FromTimestamp(time.Now())

	nintendoNotificationEventGeneral := nexproto.NewNintendoNotificationEventGeneral()

	nintendoNotificationEventGeneral.U32Param = 0
	nintendoNotificationEventGeneral.U64Param1 = 0
	nintendoNotificationEventGeneral.U64Param2 = lastOnline.Value()
	nintendoNotificationEventGeneral.StrParam = ""

	eventObject := nexproto.NewNintendoNotificationEvent()
	eventObject.Type = 10
	eventObject.SenderPID = client.PID()
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

	connectedUser := globals.ConnectedUsers[pid]

	if connectedUser != nil {
		requestPacket, _ := nex.NewPacketV0(connectedUser.Client, nil)

		requestPacket.SetVersion(0)
		requestPacket.SetSource(0xA1)
		requestPacket.SetDestination(0xAF)
		requestPacket.SetType(nex.DataPacket)
		requestPacket.SetPayload(rmcRequestBytes)

		requestPacket.AddFlag(nex.FlagNeedsAck)
		requestPacket.AddFlag(nex.FlagReliable)

		globals.NEXServer.Send(requestPacket)
	}
}
