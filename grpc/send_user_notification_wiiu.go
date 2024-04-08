package grpc

import (
	"context"

	"github.com/PretendoNetwork/friends/globals"
	pb "github.com/PretendoNetwork/grpc-go/friends"
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/constants"
	nintendo_notifications "github.com/PretendoNetwork/nex-protocols-go/v2/nintendo-notifications"
	empty "github.com/golang/protobuf/ptypes/empty"
)

// SendUserNotificationWiiU implements helloworld.SendUserNotificationWiiU
func (s *gRPCFriendsServer) SendUserNotificationWiiU(ctx context.Context, in *pb.SendUserNotificationWiiURequest) (*empty.Empty, error) {
	connectedUser := globals.ConnectedUsers[in.GetPid()]

	if connectedUser != nil {
		rmcRequest := nex.NewRMCRequest(globals.SecureEndpoint)
		rmcRequest.ProtocolID = nintendo_notifications.ProtocolID
		rmcRequest.CallID = 3810693103
		rmcRequest.MethodID = nintendo_notifications.MethodProcessNintendoNotificationEvent2
		rmcRequest.Parameters = in.GetNotificationData()

		rmcRequestBytes := rmcRequest.Bytes()

		requestPacket, _ := nex.NewPRUDPPacketV0(globals.SecureEndpoint.Server, connectedUser.Connection, nil)

		requestPacket.SetType(constants.DataPacket)
		requestPacket.AddFlag(constants.PacketFlagNeedsAck)
		requestPacket.AddFlag(constants.PacketFlagReliable)
		requestPacket.SetSourceVirtualPortStreamType(connectedUser.Connection.StreamType)
		requestPacket.SetSourceVirtualPortStreamID(globals.SecureEndpoint.StreamID)
		requestPacket.SetDestinationVirtualPortStreamType(connectedUser.Connection.StreamType)
		requestPacket.SetDestinationVirtualPortStreamID(connectedUser.Connection.StreamID)
		requestPacket.SetPayload(rmcRequestBytes)

		globals.SecureServer.Send(requestPacket)
	}

	return &empty.Empty{}, nil
}
