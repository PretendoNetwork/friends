package grpc

import (
	"context"

	"github.com/PretendoNetwork/friends/globals"
	pb "github.com/PretendoNetwork/grpc-go/friends"
	nex "github.com/PretendoNetwork/nex-go"
	nintendo_notifications "github.com/PretendoNetwork/nex-protocols-go/nintendo-notifications"
	empty "github.com/golang/protobuf/ptypes/empty"
)

// SendUserNotificationWiiU implements helloworld.SendUserNotificationWiiU
func (s *gRPCFriendsServer) SendUserNotificationWiiU(ctx context.Context, in *pb.SendUserNotificationWiiURequest) (*empty.Empty, error) {
	connectedUser := globals.ConnectedUsers[in.GetPid()]

	if connectedUser != nil {
		rmcRequest := nex.NewRMCRequest(globals.SecureServer)
		rmcRequest.ProtocolID = nintendo_notifications.ProtocolID
		rmcRequest.CallID = 3810693103
		rmcRequest.MethodID = nintendo_notifications.MethodProcessNintendoNotificationEvent2
		rmcRequest.Parameters = in.GetNotificationData()

		rmcRequestBytes := rmcRequest.Bytes()

		requestPacket, _ := nex.NewPRUDPPacketV0(connectedUser.Client, nil)

		requestPacket.SetType(nex.DataPacket)
		requestPacket.AddFlag(nex.FlagNeedsAck)
		requestPacket.AddFlag(nex.FlagReliable)
		requestPacket.SetSourceStreamType(connectedUser.Client.DestinationStreamType)
		requestPacket.SetSourcePort(connectedUser.Client.DestinationPort)
		requestPacket.SetDestinationStreamType(connectedUser.Client.SourceStreamType)
		requestPacket.SetDestinationPort(connectedUser.Client.SourcePort)
		requestPacket.SetPayload(rmcRequestBytes)

		globals.SecureServer.Send(requestPacket)
	}

	return &empty.Empty{}, nil
}
