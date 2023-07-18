package grpc

import (
	"fmt"
	"log"
	"net"
	"os"

	pb "github.com/PretendoNetwork/grpc-go/friends"
	"google.golang.org/grpc"
)

type gRPCFriendsServer struct {
	pb.UnimplementedFriendsServer
}

func StartGRPCServer() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("PN_FRIENDS_GRPC_SERVER_PORT")))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(apiKeyInterceptor),
	)

	pb.RegisterFriendsServer(server, &gRPCFriendsServer{})

	log.Printf("server listening at %v", listener.Addr())

	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
