package main

import (
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func assignNEXProtocols() {
	secureServer := nexproto.NewSecureProtocol(nexServer)
	accountManagementServer := nexproto.NewAccountManagementProtocol(nexServer)
	friendsServer := nexproto.NewFriendsProtocol(nexServer)
	friends3DSServer := nexproto.NewFriends3DSProtocol(nexServer)

	// Account Management protocol handles
	accountManagementServer.NintendoCreateAccount(nintendoCreateAccount)

	// Secure protocol handles
	secureServer.Register(register)
	secureServer.RegisterEx(registerEx)

	// Friends (WiiU) protocol handles
	friendsServer.UpdateAndGetAllInformation(updateAndGetAllInformation)
	friendsServer.AddFriendRequest(addFriendRequest)
	friendsServer.AcceptFriendRequest(acceptFriendRequest)
	friendsServer.MarkFriendRequestsAsReceived(markFriendRequestsAsReceived)
	friendsServer.UpdatePresence(updatePresenceWiiU)
	friendsServer.UpdateComment(updateCommentWiiU)
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
	friends3DSServer.UpdateComment(updateComment3DS)
}
