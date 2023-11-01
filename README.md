# ctxzap

Bind zap to context, built-in support for grpc.

## Quick start

1. Install

   ```bash
   go get github.com/MoeGolibrary/ctxzap
   ```

2. Bootstrap

    ```go
    logger := zap.NewExample()
    _ = grpc.NewServer(
        grpc.ChainUnaryInterceptor(
            grpc_zap.UnaryServerInterceptor(logger, grpc_zap.WithAccessLog(true, true)),
        ),
        grpc.ChainStreamInterceptor(
            grpc_zap.StreamServerInterceptor(logger, grpc_zap.WithEventLog(true, true)),
        ),
    )

    ```

3. Logging

    ```go
    func (s xxxServer) Ping(ctx context.Context) (xxx, error) {
        ctxzap.Info("Hello world", zap.String("id", "1"))
    }

    ```

## License

MIT
