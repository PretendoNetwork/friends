package main

import (
	"encoding/hex"
	"fmt"

	nex "github.com/PretendoNetwork/nex-go"
)

func sendTestNotification(client *nex.Client) {
	/*
		fmt.Println("Sending notification")
		///////////////////////////////////////
		// Build RMC Request for the console //
		///////////////////////////////////////

		// Setup the RMC data
		var callID uint32 = 3810693103 // hard-coded for testing
		var protocolID uint8 = 0x0E    // NotificationsProtocol
		var methodID uint32 = 0x1      // ProcessNotificationEvent

		// Create NotificationEvent
		notificationEvent := nex.NewStreamOut(nil)

		notificationEvent.WriteUInt32LE(1743126339)     // bells1998 hard-coded PID
		notificationEvent.WriteUInt32LE(0x65)           // notification type
		notificationEvent.WriteUInt64LE(49311986)       // gathering ID
		notificationEvent.WriteUInt64LE(1750087940)     // PID of user notification is being sent to
		notificationEvent.WriteString("Invite Request") // PID of user notification is being sent to

		notificationEventBytes := notificationEvent.Bytes()

		rmcRequest := nex.NewStreamOut(nil)

		rmcRequest.WriteUInt32LE(uint32(len(notificationEventBytes) + 9))
		rmcRequest.WriteUInt8(protocolID | 0x80)
		rmcRequest.WriteUInt32LE(callID)
		rmcRequest.WriteUInt32LE(methodID)
		rmcRequest.Grow(int64(len(notificationEventBytes)))
		rmcRequest.WriteBytesNext(notificationEventBytes)

		rmcRequestBytes := rmcRequest.Bytes()

		fmt.Println(hex.EncodeToString(rmcRequestBytes))

		requestPacket, _ := nex.NewPacketV0(client, nil)

		requestPacket.SetVersion(0)
		requestPacket.SetSource(0xA1)
		requestPacket.SetDestination(0xAF)
		requestPacket.SetType(nex.DataPacket)
		requestPacket.SetPayload(rmcRequestBytes)

		requestPacket.AddFlag(nex.FlagNeedsAck)
		requestPacket.AddFlag(nex.FlagReliable)

		nexServer.Send(requestPacket)
	*/

	fmt.Println("Sending notification")
	///////////////////////////////////////
	// Build RMC Request for the console //
	///////////////////////////////////////

	// Setup the RMC data
	var callID uint32 = 3810693103 // hard-coded for testing
	var protocolID uint8 = 0x64    // NotificationsProtocol
	var methodID uint32 = 0x1      // ProcessNotificationEvent

	// Create NintendoNotificationEventGeneral DataHolder

	datetime := nex.NewDateTime(0)
	nintendoNotificationEventGeneral := nex.NewStreamOut(nil)

	nintendoNotificationEventGeneral.WriteUInt32LE(1750087940)     // PID of user notification is being sent to
	nintendoNotificationEventGeneral.WriteUInt64LE(49311986)       // gathering ID
	nintendoNotificationEventGeneral.WriteUInt64LE(datetime.Now()) // DateTime of notification
	nintendoNotificationEventGeneral.WriteString("Invite Request") // Notification message

	nintendoNotificationEventGeneralBytes := nintendoNotificationEventGeneral.Bytes()

	// Create NintendoNotificationEvent

	nintendoNotificationEvent := nex.NewStreamOut(nil)

	nintendoNotificationEvent.WriteUInt32LE(0x65)       // Notification type
	nintendoNotificationEvent.WriteUInt32LE(1743126339) // bells1998 hard-coded PID

	// Write NintendoNotificationEventGeneral DataHolder

	nintendoNotificationEvent.WriteString("NintendoNotificationEventGeneral")
	nintendoNotificationEvent.WriteUInt32LE(uint32(len(nintendoNotificationEventGeneralBytes) + 8))
	nintendoNotificationEvent.WriteBuffer(nintendoNotificationEventGeneralBytes)

	nintendoNotificationEventBytes := nintendoNotificationEvent.Bytes()

	// Setup RMC Request

	rmcRequest := nex.NewStreamOut(nil)

	rmcRequest.WriteUInt32LE(uint32(len(nintendoNotificationEventBytes) + 9))
	rmcRequest.WriteUInt8(protocolID | 0x80)
	rmcRequest.WriteUInt32LE(callID)
	rmcRequest.WriteUInt32LE(methodID)
	rmcRequest.Grow(int64(len(nintendoNotificationEventBytes)))
	rmcRequest.WriteBytesNext(nintendoNotificationEventBytes)

	rmcRequestBytes := rmcRequest.Bytes()

	fmt.Println(hex.EncodeToString(rmcRequestBytes))

	requestPacket, _ := nex.NewPacketV0(client, nil)

	requestPacket.SetVersion(0)
	requestPacket.SetSource(0xA1)
	requestPacket.SetDestination(0xAF)
	requestPacket.SetType(nex.DataPacket)
	requestPacket.SetPayload(rmcRequestBytes)

	requestPacket.AddFlag(nex.FlagNeedsAck)
	requestPacket.AddFlag(nex.FlagReliable)

	nexServer.Send(requestPacket)
}
