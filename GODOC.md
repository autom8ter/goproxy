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

#### func  New

```go
func New(config *ProxyConfig) *GoProxy
```
New registers a new reverseproxy for each provided ProxyConfig

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

#### func (*GoProxy) ModifyRequests

```go
func (g *GoProxy) ModifyRequests(middleware middleware.RequestWare)
```
ModifyResponses takes a Request Middleware function, traverses each registered
reverse proxy, and modifies the http request it sends to its target prior to
sending

#### func (*GoProxy) ModifyResponses

```go
func (g *GoProxy) ModifyResponses(middleware middleware.ResponseWare)
```
ModifyResponses takes a Response Middleware function, traverses each registered
reverse proxy, and modifies the http response it sends to the client

#### func (*GoProxy) ModifyRouter

```go
func (g *GoProxy) ModifyRouter(middleware middleware.RouterWare)
```
ModifyRouter takes a router middleware function and wraps the proxies router

#### func (*GoProxy) ModifyTransport

```go
func (g *GoProxy) ModifyTransport(middleware middleware.TransportWare)
```
ModifyResponses takes a Transport Middleware function, traverses each registered
reverse proxy, and modifies the http roundtripper it uses

#### func (*GoProxy) Proxies

```go
func (g *GoProxy) Proxies() map[string]*httputil.ReverseProxy
```
Proxies returns all registered reverse proxies as a map of prefix:reverse proxy

#### func (*GoProxy) WalkPaths

```go
func (g *GoProxy) WalkPaths(fns ...mux.WalkFunc) error
```
WalkPaths walks registered mux paths and modifies them

#### type ProxyConfig

```go
type ProxyConfig struct {
	Configs []*Config
}
```

ProxyConfig configures the entire reverse proxy
