package main

import (
	"fmt"

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
	nexServer.SetAccessKey("ridfebb9")

	nexServer.On("Data", func(packet *nex.PacketV0) {
		request := packet.RMCRequest()

		fmt.Println("==Friends - Secure==")
		fmt.Printf("Protocol ID: %#v\n", request.ProtocolID())
		fmt.Printf("Method ID: %#v\n", request.MethodID())
		fmt.Println("====================")
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
	friendsServer.MarkFriendRequestsAsReceived(markFriendRequestsAsReceived)
	friendsServer.UpdatePreference(updatePreferenceWiiU)
	friendsServer.GetBasicInfo(getBasicInfo)
	friendsServer.CheckSettingStatus(checkSettingStatus)
	friendsServer.GetRequestBlockSettings(getRequestBlockSettings)

	// Friends (3DS) protocol handles
	friends3DSServer.UpdateProfile(updateProfile)
	friends3DSServer.UpdateMii(updateMii)
	friends3DSServer.UpdatePreference(updatePreference3DS)
	friends3DSServer.SyncFriend(syncFriend)
	friends3DSServer.UpdatePresence(updatePresence)
	friends3DSServer.UpdateFavoriteGameKey(updateFavoriteGameKey)
	friends3DSServer.UpdateComment(updateComment)

	nexServer.Listen(":60001")
}
