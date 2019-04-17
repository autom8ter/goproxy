package grpc

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"net/http"
	"strings"
)

func MustDial(ctx context.Context, addr string, opts ...grpc.DialOption) *grpc.ClientConn {
	conn, err := grpc.DialContext(ctx, "", opts...)
	if err != nil {
		panic(err.Error())
	}
	return conn
}

func HandleGRPC(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	})
}

func Gateway() *runtime.ServeMux {
	return runtime.NewServeMux()
}
