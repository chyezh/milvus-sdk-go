package client

import (
	"context"
	"fmt"

	"github.com/milvus-io/milvus-sdk-go/v2/internal/utils/crypto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// authenticationInterceptor appends credential into context metadata
func authenticationInterceptor(ctx context.Context, username, password string) context.Context {
	value := crypto.Base64Encode(fmt.Sprintf("%s:%s", username, password))
	return metadata.AppendToOutgoingContext(ctx, "authorization", value)
}

// CreateAuthenticationUnaryInterceptor creates a unary interceptor for authentication
func createAuthenticationUnaryInterceptor(username, password string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = authenticationInterceptor(ctx, username, password)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// createAuthenticationStreamInterceptor creates a stream interceptor for authentication
func createAuthenticationStreamInterceptor(username, password string) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		ctx = authenticationInterceptor(ctx, username, password)
		return streamer(ctx, desc, cc, method, opts...)
	}
}

// databaseNameInterceptor appends the dbName into metadata.
func databaseNameInterceptor(ctx context.Context, dbName string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "dbname", dbName)
}

// createDatabaseNameInterceptor creates a unary interceptor for db name.
func createDatabaseNameUnaryInterceptor(dbName string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = databaseNameInterceptor(ctx, dbName)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// createDatabaseNameStreamInterceptor creates a unary interceptor for db name.
func createDatabaseNameStreamInterceptor(dbName string) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		ctx = databaseNameInterceptor(ctx, dbName)
		return streamer(ctx, desc, cc, method, opts...)
	}
}

// cloudMetaInterceptor appends the cloud meta into metadata.
func cloudMetaInterceptor(ctx context.Context, meta *CloudConfig) context.Context {
	ctx = metadata.AppendToOutgoingContext(ctx, "cloud-meta-api-key", meta.APIKey)
	return metadata.AppendToOutgoingContext(ctx, "cloud-meta-cluster", meta.ClusterName)
}

// createCloudMetaInterceptor creates a unary interceptor for db name.
func createCloudMetaUnaryInterceptor(meta *CloudConfig) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = cloudMetaInterceptor(ctx, meta)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// createCloudMetaStreamInterceptor creates a unary interceptor for db name.
func createCloudMetaStreamInterceptor(meta *CloudConfig) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		ctx = cloudMetaInterceptor(ctx, meta)
		return streamer(ctx, desc, cc, method, opts...)
	}
}
