package grpc

import (
	"context"

	"github.com/PretendoNetwork/friends-secure/globals"
	pb "github.com/PretendoNetwork/grpc-go/friends"
	nex "github.com/PretendoNetwork/nex-go"
	nintendo_notifications "github.com/PretendoNetwork/nex-protocols-go/nintendo-notifications"
	empty "github.com/golang/protobuf/ptypes/empty"
)

// SendUserNotificationWiiU implements helloworld.SendUserNotificationWiiU
func (s *gRPCFriendsServer) SendUserNotificationWiiU(ctx context.Context, in *pb.SendUserNotificationWiiURequest) (*empty.Empty, error) {
	connectedUser := globals.ConnectedUsers[in.GetPid()]

	if connectedUser != nil {
		rmcRequest := nex.NewRMCRequest()
		rmcRequest.SetProtocolID(nintendo_notifications.ProtocolID)
		rmcRequest.SetCallID(3810693103)
		rmcRequest.SetMethodID(nintendo_notifications.MethodProcessNintendoNotificationEvent2)
		rmcRequest.SetParameters(in.GetNotificationData())

		rmcRequestBytes := rmcRequest.Bytes()

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

	return &empty.Empty{}, nil
}
