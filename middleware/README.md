# middleware
--
    import "github.com/autom8ter/goproxy/middleware"


## Usage

#### type HandlerWare

```go
type HandlerWare func(h http.Handler) http.Handler
```

RouterWare is a function used to modify the mux

#### func  WithCORS

```go
func WithCORS(options cors.Options) HandlerWare
```

#### func  WithJWT

```go
func WithJWT(secret string) HandlerWare
```
