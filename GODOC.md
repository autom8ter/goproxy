# goproxy
--
    import "github.com/autom8ter/goproxy"


## Usage

#### type GoProxy

```go
type GoProxy struct {
	*http.ServeMux
}
```

GoProxy is an API Gateway/Reverse Proxy and http.ServeMux/http.Handler

#### func  NewGoProxy

```go
func NewGoProxy(configs ...*ProxyConfig) *GoProxy
```
NewGoProxy registers a new reverseproxy for each provided ProxyConfig

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
ListenAndServe starts an http server on the given address with all of the
registered reverse proxies

#### func (*GoProxy) ModifyRequests

```go
func (g *GoProxy) ModifyRequests(middleware RequestMiddleware)
```
ModifyResponses takes a Request Middleware function, traverses each registered
reverse proxy, and modifies the http request it sends to its target prior to
sending

#### func (*GoProxy) ModifyResponses

```go
func (g *GoProxy) ModifyResponses(middleware ResponseMiddleware)
```
ModifyResponses takes a Response Middleware function, traverses each registered
reverse proxy, and modifies the http response it sends to the client

#### func (*GoProxy) ModifyTransport

```go
func (g *GoProxy) ModifyTransport(middleware TransportMiddleware)
```
ModifyResponses takes a Transport Middleware function, traverses each registered
reverse proxy, and modifies the http roundtripper it uses

#### func (*GoProxy) Proxies

```go
func (g *GoProxy) Proxies() map[string]*httputil.ReverseProxy
```
Proxies returns all registered reverse proxies as a map of prefix:reverse proxy

#### func (*GoProxy) WithPprof

```go
func (g *GoProxy) WithPprof() *GoProxy
```

#### type ProxyConfig

```go
type ProxyConfig struct {
	PathPrefix string
	TargetUrl  string
	Username   string
	Password   string
	Headers    map[string]string
}
```

ProxyConfig is used to configure GoProxies reverse proxies

#### type RequestMiddleware

```go
type RequestMiddleware func(func(req *http.Request)) func(req *http.Request)
```

RequestMiddleware is a function used to modify the incoming request of a reverse
proxy from a client

#### type ResponseMiddleware

```go
type ResponseMiddleware func(func(response *http.Response) error) func(response *http.Response) error
```

ResponseMiddleware is a function used to modify the response of a reverse proxy

#### type TransportMiddleware

```go
type TransportMiddleware func(tripper http.RoundTripper) http.RoundTripper
```

TransportMiddleware is a function used to modify the http RoundTripper that is
used by a reverse proxy. The default RoundTripper is initially
http.DefaultTransport
