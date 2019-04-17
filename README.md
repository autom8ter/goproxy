# goproxy
--
    import "github.com/autom8ter/goproxy"


## Usage

#### type Config

```go
type Config struct {
	PathPrefix string `validate:"required"`
	TargetUrl  string `validate:"required"`
	Username   string
	Password   string
	Headers    map[string]string
	FormValues map[string]string
}
```

Config is used to configure a reverse proxy handler(one route)

#### type GoProxy

```go
type GoProxy struct {
	*mux.Router
}
```

GoProxy is an API Gateway/Reverse Proxy and http.ServeMux/http.Handler

#### func  NewGoProxy

```go
func NewGoProxy(configs ...*Config) *GoProxy
```
NewGoProxy registers a new reverseproxy handler for each provided config with
the specified path prefix

#### func  NewSecureGoProxy

```go
func NewSecureGoProxy(secret string, opts cors.Options, configs ...*Config) *GoProxy
```
NewSecureGoProxy registers a new secure reverseproxy for each provided configs.
It is the same as New, except with CORS options and a JWT middleware that checks
for a signed bearer token

#### func (*GoProxy) AsHandlerFunc

```go
func (g *GoProxy) AsHandlerFunc() http.HandlerFunc
```
AsHandlerFunc converts a GoProxy to an http.HandlerFunc for convenience

#### func (*GoProxy) GetProxy

```go
func (g *GoProxy) GetProxy(prefix string) *httputil.ReverseProxy
```
GetProxy returns the reverse proxy with the registered prefix

#### func (*GoProxy) ListenAndServe

```go
func (g *GoProxy) ListenAndServe(addr string) error
```
ListenAndServe starts the GoProxy server on the specified address

#### func (*GoProxy) Middleware

```go
func (g *GoProxy) Middleware(middlewares ...mux.MiddlewareFunc)
```
Middleware wraps Goproxy with the provided middlewares

#### func (*GoProxy) Proxies

```go
func (g *GoProxy) Proxies() map[string]*httputil.ReverseProxy
```
Proxies returns all registered reverse proxies as a map of prefix:reverse proxy

#### func (*GoProxy) RequestWare

```go
func (g *GoProxy) RequestWare(middleware middleware.RequestWare)
```
ModifyResponses takes a Request Middleware function, traverses each registered
reverse proxy, and modifies the http request it sends to its target prior to
sending

#### func (*GoProxy) ResponseWare

```go
func (g *GoProxy) ResponseWare(middleware middleware.ResponseWare)
```
ModifyResponses takes a Response Middleware function, traverses each registered
reverse proxy, and modifies the http response it sends to the client

#### func (*GoProxy) TransportWare

```go
func (g *GoProxy) TransportWare(middleware middleware.TransportWare)
```
ModifyResponses takes a Transport Middleware function, traverses each registered
reverse proxy, and modifies the http roundtripper it uses

#### func (*GoProxy) WalkPaths

```go
func (g *GoProxy) WalkPaths(walkfuncs ...mux.WalkFunc) error
```
WalkPaths walks registered mux paths

#### func (*GoProxy) WithMetrics

```go
func (g *GoProxy) WithMetrics()
```
Registers prometheus metrics for: in_flight_requests, requests_total,
request_duration_seconds, response_size_bytes,

#### func (*GoProxy) WithPprof

```go
func (g *GoProxy) WithPprof()
```
registers all pprof handlers: /debug/pprof/, /debug/pprof/cmdline,
/debug/pprof/profile, /debug/pprof/symbol, /debug/pprof/trace
