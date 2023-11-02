# ctxzap

Bind zap to context, built-in support for grpc.

## Quick start

1. Install

   ```bash
   go get github.com/acrazing/ctxzap
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

    // In any context based method called by grpc servers
    _ = func(ctx context.Context) {
        // add a field to the context, this field will always be printed in the next
        // log calls with the context.
        ctxzap.AddFields(ctx, zap.String("data", "1"))

        // replace an existing field or append it to the context
        ctxzap.ReplaceField(ctx, zap.String("data", "2"))

        // print a log
        ctxzap.Info(ctx, "call with data", zap.String("value", "2"))
    }

    ```

3. Logging

    ```go
    func (s xxxServer) Ping(ctx context.Context) (xxx, error) {
        ctxzap.Info("Hello world", zap.String("id", "1"))
    }

    ```

## License

MIT
