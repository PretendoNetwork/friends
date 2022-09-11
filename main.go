package main

import (
	"sync"
	"time"

	database_wiiu "github.com/PretendoNetwork/friends-secure/database/wiiu"
	"github.com/PretendoNetwork/friends-secure/globals"
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

var wg sync.WaitGroup

func main() {
	wg.Add(2)

	go startGRPCServer()
	go startNEXServer()

	wg.Wait()
}

// Maybe this function should go in a different file?
func sendUserWentOfflineWiiUNotifications(client *nex.Client) {
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

	friendList := database_wiiu.GetUserFriendList(client.PID())

	for i := 0; i < len(friendList); i++ {
		friendPID := friendList[i].NNAInfo.PrincipalBasicInfo.PID
		connectedUser := globals.ConnectedUsers[friendPID]

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
}
