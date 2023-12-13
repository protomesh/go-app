package app

import (
	"context"

	"google.golang.org/grpc"
)

func GrpcPathPrefixUnaryInterceptor(prefix string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		modifiedMethod := prefix + method
		return invoker(ctx, modifiedMethod, req, reply, cc, opts...)
	}
}

func GrpcPathPrefixStreamInterceptor(prefix string) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		modifiedMethod := prefix + method
		return streamer(ctx, desc, cc, modifiedMethod, opts...)
	}
}
