package main

import (
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func assignNEXProtocols() {
	secureServer := nexproto.NewSecureProtocol(nexServer)
	accountManagementServer := nexproto.NewAccountManagementProtocol(nexServer)
	friendsWiiUServer := nexproto.NewFriendsWiiUProtocol(nexServer)
	friends3DSServer := nexproto.NewFriends3DSProtocol(nexServer)

	// Account Management protocol handles
	accountManagementServer.NintendoCreateAccount(nintendoCreateAccount)

	// Secure protocol handles
	secureServer.Register(register)
	secureServer.RegisterEx(registerEx)

	// Friends (WiiU) protocol handles
	friendsWiiUServer.UpdateAndGetAllInformation(updateAndGetAllInformation)
	friendsWiiUServer.AddFriendRequest(addFriendRequest)
	friendsWiiUServer.AcceptFriendRequest(acceptFriendRequest)
	friendsWiiUServer.MarkFriendRequestsAsReceived(markFriendRequestsAsReceived)
	friendsWiiUServer.UpdatePresence(updatePresenceWiiU)
	friendsWiiUServer.UpdateComment(updateCommentWiiU)
	friendsWiiUServer.UpdatePreference(updatePreferenceWiiU)
	friendsWiiUServer.GetBasicInfo(getBasicInfo)
	friendsWiiUServer.DeletePersistentNotification(deletePersistentNotification)
	friendsWiiUServer.CheckSettingStatus(checkSettingStatus)
	friendsWiiUServer.GetRequestBlockSettings(getRequestBlockSettings)

	// Friends (3DS) protocol handles
	friends3DSServer.UpdateProfile(updateProfile)
	friends3DSServer.UpdateMii(updateMii)
	friends3DSServer.UpdatePreference(updatePreference3DS)
	friends3DSServer.SyncFriend(syncFriend)
	friends3DSServer.UpdatePresence(updatePresence3DS)
	friends3DSServer.UpdateFavoriteGameKey(updateFavoriteGameKey)
	friends3DSServer.UpdateComment(updateComment3DS)
}
