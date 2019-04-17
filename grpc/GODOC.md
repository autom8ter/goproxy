# grpc
--
    import "github.com/autom8ter/goproxy/grpc"


## Usage

#### func  Gateway

```go
func Gateway() *runtime.ServeMux
```

#### func  HandleGRPC

```go
func HandleGRPC(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler
```

#### func  MustDial

```go
func MustDial(ctx context.Context, addr string, opts ...grpc.DialOption) *grpc.ClientConn
```
