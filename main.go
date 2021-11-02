package main

import (
	"fmt"
	"time"

	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

var nexServer *nex.Server
var secureServer *nexproto.SecureProtocol

func main() {
	nexServer = nex.NewServer()
	nexServer.SetPrudpVersion(0)
	nexServer.SetSignatureVersion(1)
	nexServer.SetKerberosKeySize(16)
	nexServer.SetPingTimeout(20) // Maybe too long?
	nexServer.SetAccessKey("ridfebb9")

	nexServer.On("Data", func(packet *nex.PacketV0) {
		request := packet.RMCRequest()

		fmt.Println("==Friends - Secure==")
		fmt.Printf("Protocol ID: %#v\n", request.ProtocolID())
		fmt.Printf("Method ID: %#v\n", request.MethodID())
		fmt.Println("====================")
	})

	nexServer.On("Kick", func(packet *nex.PacketV0) {
		pid := packet.Sender().PID()
		delete(connectedUsers, pid)

		lastOnline := nex.NewDateTime(0)
		lastOnline.FromTimestamp(time.Now())

		updateUserLastOnlineTime(pid, lastOnline)
		sendUserWentOfflineWiiUNotifications(packet.Sender())

		fmt.Println("Leaving")
	})

	nexServer.On("Ping", func(packet *nex.PacketV0) {
		fmt.Print("Pinged. Is ACK: ")
		fmt.Println(packet.HasFlag(nex.FlagAck))
	})

	secureServer = nexproto.NewSecureProtocol(nexServer)
	accountManagementServer := nexproto.NewAccountManagementProtocol(nexServer)
	friendsServer := nexproto.NewFriendsProtocol(nexServer)
	friends3DSServer := nexproto.NewFriends3DSProtocol(nexServer)

	// Handle PRUDP CONNECT packet (not an RMC method)
	nexServer.On("Connect", connect)

	// Account Management protocol handles
	accountManagementServer.NintendoCreateAccount(nintendoCreateAccount)
	accountManagementServer.NintendoCreateAccount3DS(nintendoCreateAccount3DS)

	// Secure protocol handles
	secureServer.Register(register)
	secureServer.RegisterEx(registerEx)

	// Friends (WiiU) protocol handles
	friendsServer.UpdateAndGetAllInformation(updateAndGetAllInformation)
	friendsServer.AddFriendRequest(addFriendRequest)
	friendsServer.AcceptFriendRequest(acceptFriendRequest)
	friendsServer.MarkFriendRequestsAsReceived(markFriendRequestsAsReceived)
	friendsServer.UpdatePresence(updatePresenceWiiU)
	friendsServer.UpdatePreference(updatePreferenceWiiU)
	friendsServer.GetBasicInfo(getBasicInfo)
	friendsServer.DeletePersistentNotification(deletePersistentNotification)
	friendsServer.CheckSettingStatus(checkSettingStatus)
	friendsServer.GetRequestBlockSettings(getRequestBlockSettings)

	// Friends (3DS) protocol handles
	friends3DSServer.UpdateProfile(updateProfile)
	friends3DSServer.UpdateMii(updateMii)
	friends3DSServer.UpdatePreference(updatePreference3DS)
	friends3DSServer.SyncFriend(syncFriend)
	friends3DSServer.UpdatePresence(updatePresence3DS)
	friends3DSServer.UpdateFavoriteGameKey(updateFavoriteGameKey)
	friends3DSServer.UpdateComment(updateComment)

	nexServer.Listen(":60001")
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
	eventObject.DataHolder.Name = "NintendoNotificationEventGeneral"
	eventObject.DataHolder.Object = nintendoNotificationEventGeneral

	stream := nex.NewStreamOut(nexServer)
	stream.WriteStructure(eventObject)

	rmcRequest := nex.NewRMCRequest()
	rmcRequest.SetProtocolID(nexproto.NintendoNotificationsProtocolID)
	rmcRequest.SetCallID(3810693103)
	rmcRequest.SetMethodID(nexproto.NintendoNotificationsMethodProcessNintendoNotificationEvent1)
	rmcRequest.SetParameters(stream.Bytes())

	rmcRequestBytes := rmcRequest.Bytes()

	friendList := getUserFriendList(client.PID())

	for i := 0; i < len(friendList); i++ {
		friendPID := friendList[i].NNAInfo.PrincipalBasicInfo.PID
		connectedUser := connectedUsers[friendPID]

		if connectedUser != nil {
			requestPacket, _ := nex.NewPacketV0(connectedUser.Client, nil)

			requestPacket.SetVersion(0)
			requestPacket.SetSource(0xA1)
			requestPacket.SetDestination(0xAF)
			requestPacket.SetType(nex.DataPacket)
			requestPacket.SetPayload(rmcRequestBytes)

			requestPacket.AddFlag(nex.FlagNeedsAck)
			requestPacket.AddFlag(nex.FlagReliable)

			nexServer.Send(requestPacket)
		}
	}
}
