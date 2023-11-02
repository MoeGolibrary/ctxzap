# ctxzap

Bind zap to context, built-in support for grpc.

## Install

   ```bash
   go get github.com/acrazing/ctxzap
   ```

## Quick start

```go
package main

import (
    "context"
    "github.com/acrazing/ctxzap"
    "github.com/acrazing/ctxzap/grpc_zap"
    "go.uber.org/zap"
    "google.golang.org/grpc"
)

func main() {
	// get your logger
    logger := zap.NewExample()

	// inject interceptors
    _ = grpc.NewServer(
        grpc.ChainUnaryInterceptor(
            grpc_zap.UnaryServerInterceptor(logger, grpc_zap.WithAccessLog(true, true)),
        ),
        grpc.ChainStreamInterceptor(
            grpc_zap.StreamServerInterceptor(logger, grpc_zap.WithEventLog(true, true)),
        ),
    )

	// print log
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
}
```

## License

MIT
