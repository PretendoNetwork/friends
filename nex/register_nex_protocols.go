package nex

import (
	"github.com/PretendoNetwork/friends-secure/globals"
	nex_account_management "github.com/PretendoNetwork/friends-secure/nex/account-management"
	nex_friends_3ds "github.com/PretendoNetwork/friends-secure/nex/friends-3ds"
	nex_friends_wiiu "github.com/PretendoNetwork/friends-secure/nex/friends-wiiu"
	nex_secure_connection "github.com/PretendoNetwork/friends-secure/nex/secure-connection"
	account_management "github.com/PretendoNetwork/nex-protocols-go/account-management"
	friends_3ds "github.com/PretendoNetwork/nex-protocols-go/friends/3ds"
	friends_wiiu "github.com/PretendoNetwork/nex-protocols-go/friends/wiiu"
	secure_connection "github.com/PretendoNetwork/nex-protocols-go/secure-connection"
)

func registerNEXProtocols() {
	secureConnectionProtocol := secure_connection.NewSecureConnectionProtocol(globals.NEXServer)
	accountManagementProtocol := account_management.NewAccountManagementProtocol(globals.NEXServer)
	friendsWiiUProtocol := friends_wiiu.NewFriendsWiiUProtocol(globals.NEXServer)
	friends3DSProtocol := friends_3ds.NewFriends3DSProtocol(globals.NEXServer)

	// Account Management protocol handles
	accountManagementProtocol.NintendoCreateAccount(nex_account_management.NintendoCreateAccount)

	// Secure protocol handles
	secureConnectionProtocol.Register(nex_secure_connection.Register)
	secureConnectionProtocol.RegisterEx(nex_secure_connection.RegisterEx)

	// Friends (WiiU) protocol handles
	friendsWiiUProtocol.UpdateAndGetAllInformation(nex_friends_wiiu.UpdateAndGetAllInformation)
	friendsWiiUProtocol.AddFriendRequest(nex_friends_wiiu.AddFriendRequest)
	friendsWiiUProtocol.RemoveFriend(nex_friends_wiiu.RemoveFriend)
	friendsWiiUProtocol.CancelFriendRequest(nex_friends_wiiu.CancelFriendRequest)
	friendsWiiUProtocol.AcceptFriendRequest(nex_friends_wiiu.AcceptFriendRequest)
	friendsWiiUProtocol.DeleteFriendRequest(nex_friends_wiiu.DeleteFriendRequest)
	friendsWiiUProtocol.DenyFriendRequest(nex_friends_wiiu.DenyFriendRequest)
	friendsWiiUProtocol.MarkFriendRequestsAsReceived(nex_friends_wiiu.MarkFriendRequestsAsReceived)
	friendsWiiUProtocol.AddBlackList(nex_friends_wiiu.AddBlacklist)
	friendsWiiUProtocol.RemoveBlackList(nex_friends_wiiu.RemoveBlacklist)
	friendsWiiUProtocol.UpdatePresence(nex_friends_wiiu.UpdatePresence)
	friendsWiiUProtocol.UpdateComment(nex_friends_wiiu.UpdateComment)
	friendsWiiUProtocol.UpdatePreference(nex_friends_wiiu.UpdatePreference)
	friendsWiiUProtocol.GetBasicInfo(nex_friends_wiiu.GetBasicInfo)
	friendsWiiUProtocol.DeletePersistentNotification(nex_friends_wiiu.DeletePersistentNotification)
	friendsWiiUProtocol.CheckSettingStatus(nex_friends_wiiu.CheckSettingStatus)
	friendsWiiUProtocol.GetRequestBlockSettings(nex_friends_wiiu.GetRequestBlockSettings)

	// Friends (3DS) protocol handles
	friends3DSProtocol.UpdateProfile(nex_friends_3ds.UpdateProfile)
	friends3DSProtocol.UpdateMii(nex_friends_3ds.UpdateMii)
	friends3DSProtocol.UpdatePreference(nex_friends_3ds.UpdatePreference)
	friends3DSProtocol.SyncFriend(nex_friends_3ds.SyncFriend)
	friends3DSProtocol.UpdatePresence(nex_friends_3ds.UpdatePresence)
	friends3DSProtocol.UpdateFavoriteGameKey(nex_friends_3ds.UpdateFavoriteGameKey)
	friends3DSProtocol.UpdateComment(nex_friends_3ds.UpdateComment)
	friends3DSProtocol.AddFriendByPrincipalID(nex_friends_3ds.AddFriendshipByPrincipalID)
	friends3DSProtocol.GetFriendPersistentInfo(nex_friends_3ds.GetFriendPersistentInfo)
	friends3DSProtocol.GetFriendMii(nex_friends_3ds.GetFriendMii)
	friends3DSProtocol.GetFriendPresence(nex_friends_3ds.GetFriendPresence)
	friends3DSProtocol.RemoveFriendByPrincipalID(nex_friends_3ds.RemoveFriendByPrincipalID)
	friends3DSProtocol.RemoveFriendByLocalFriendCode(nex_friends_3ds.RemoveFriendByLocalFriendCode)
	friends3DSProtocol.GetPrincipalIDByLocalFriendCode(nex_friends_3ds.GetPrincipalIDByLocalFriendCode)
	friends3DSProtocol.GetAllFriends(nex_friends_3ds.GetAllFriends)
}
