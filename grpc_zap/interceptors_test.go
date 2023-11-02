package grpc_zap_test

import (
	"context"
	"fmt"
	"github.com/acrazing/ctxzap/grpc_zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"net"
	"os"
	"testing"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/testing/testpb"

	"github.com/stretchr/testify/require"
)

func ExampleStreamServerInterceptor() {
	logger := zap.NewExample()
	_ = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_zap.UnaryServerInterceptor(logger, grpc_zap.WithAccessLog(true, true)),
		),
		grpc.ChainStreamInterceptor(
			grpc_zap.StreamServerInterceptor(logger, grpc_zap.WithEventLog(true, true)),
		),
	)
}

func TestStreamServerInterceptor(t *testing.T) {
	stopped := make(chan error)
	writeSyncer := zapcore.AddSync(os.Stdout)
	logLevel := zap.DebugLevel
	logEncoderConfig := zap.NewProductionEncoderConfig()
	logEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logEncoderConfig.EncodeDuration = zapcore.MillisDurationEncoder

	core := zapcore.NewCore(zapcore.NewJSONEncoder(logEncoderConfig), writeSyncer, logLevel)

	logger := zap.New(core)

	serverListener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_zap.UnaryServerInterceptor(logger, grpc_zap.WithAccessLog(true, true), grpc_zap.WithMetadataFields("traceID"), grpc_zap.WithErrorFieldsFunc(func(ctx context.Context, err error) []zap.Field {
				return []zap.Field{
					zap.Stack("stack"),
				}
			})),
		),
		grpc.ChainStreamInterceptor(
			grpc_zap.StreamServerInterceptor(logger, grpc_zap.WithEventLog(true, true), grpc_zap.WithMetadataFields("deviceID")),
		),
	)
	testpb.RegisterTestServiceServer(server, &testpb.TestPingService{})

	go func() {
		defer close(stopped)
		stopped <- server.Serve(serverListener)
	}()
	defer func() {
		server.Stop()
		<-stopped
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// This is the point where we hook up the interceptor.
	clientConn, err := grpc.DialContext(
		ctx,
		serverListener.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	require.NoError(t, err, "must not error on client Dial")

	testClient := testpb.NewTestServiceClient(clientConn)
	select {
	case err := <-stopped:
		t.Fatal("gRPC server stopped prematurely", err)
	default:
	}

	r, err := testClient.PingEmpty(context.Background(), &testpb.PingEmptyRequest{})
	require.NoError(t, err)
	require.NotNil(t, r)

	r2, err := testClient.Ping(context.Background(), &testpb.PingRequest{Value: "24"})
	require.NoError(t, err)
	require.Equal(t, "24", r2.Value)
	require.Equal(t, int32(0), r2.Counter)

	_, err = testClient.PingError(context.Background(), &testpb.PingErrorRequest{
		ErrorCodeReturned: uint32(codes.Internal),
		Value:             "24",
	})
	require.Error(t, err)
	require.Equal(t, codes.Internal, status.Code(err))

	l, err := testClient.PingList(context.Background(), &testpb.PingListRequest{Value: "24"})
	require.NoError(t, err)
	for i := 0; i < testpb.ListResponseCount; i++ {
		r, err := l.Recv()
		require.NoError(t, err)
		require.Equal(t, "24", r.Value)
		require.Equal(t, int32(i), r.Counter)
	}

	s, err := testClient.PingStream(context.Background())
	require.NoError(t, err)
	for i := 0; i < testpb.ListResponseCount; i++ {
		require.NoError(t, s.Send(&testpb.PingStreamRequest{Value: fmt.Sprintf("%v", i)}))

		r, err := s.Recv()
		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf("%v", i), r.Value)
		require.Equal(t, int32(i), r.Counter)
	}

	select {
	case err := <-stopped:
		t.Fatal("gRPC server stopped prematurely", err)
	default:
	}
}
