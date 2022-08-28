package main

import (
	"context"
	"log"
	"net"

	pb "github.com/PretendoNetwork/grpc-go/friends"
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
	empty "github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

type gRPCFriendsServer struct {
	pb.UnimplementedFriendsServer
}

// SendUserNotificationWiiU implements helloworld.SendUserNotificationWiiU
func (s *gRPCFriendsServer) SendUserNotificationWiiU(ctx context.Context, in *pb.SendUserNotificationWiiURequest) (*empty.Empty, error) {
	connectedUser := connectedUsers[in.Pid]

	if connectedUser != nil {
		rmcRequest := nex.NewRMCRequest()
		rmcRequest.SetProtocolID(nexproto.NintendoNotificationsProtocolID)
		rmcRequest.SetCallID(3810693103)
		rmcRequest.SetMethodID(nexproto.NintendoNotificationsMethodProcessNintendoNotificationEvent1)
		rmcRequest.SetParameters(in.NotificationData)

		rmcRequestBytes := rmcRequest.Bytes()

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

	return &empty.Empty{}, nil
}

func gRPCStart() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()

	pb.RegisterFriendsServer(server, &gRPCFriendsServer{})

	log.Printf("server listening at %v", listener.Addr())

	server.Serve(listener)
}
