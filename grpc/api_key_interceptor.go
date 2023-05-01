package grpc

import (
	"context"
	"errors"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func apiKeyInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)

	if ok {
		apiKeyHeader := md.Get("X-API-Key")

		if len(apiKeyHeader) == 0 || apiKeyHeader[0] != os.Getenv("PN_FRIENDS_CONFIG_GRPC_API_KEY") {
			return nil, errors.New("Missing or invalid API key")
		}
	}

	h, err := handler(ctx, req)

	return h, err
}
