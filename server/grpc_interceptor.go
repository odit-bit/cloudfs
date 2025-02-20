package server

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type interceptor struct {
}

func NewInterceptor() *interceptor {
	i := &interceptor{}
	return i
}

func (i *interceptor) InterceptStream(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	md, ok := metadata.FromIncomingContext(ss.Context())
	if !ok {
		return status.Errorf(codes.InvalidArgument, "missing metadata")
	}
	if !(i.valid(md["authorization"])) {
		return status.Errorf(codes.Unauthenticated, "invalid token")
	}
	return handler(srv, ss)
}
func (i *interceptor) Intercept(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "missing metadata")
	}
	if !(i.valid(md["authorization"])) {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}
	return handler(ctx, req)
}

func (i *interceptor) valid(authz []string) bool {
	if len(authz) < 1 {
		return false
	}

	token := strings.TrimPrefix(authz[0], "Bearer ")
	return token != ""
}
