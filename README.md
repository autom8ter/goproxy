# goproxy
--
    import "github.com/autom8ter/goproxy"


## Usage

#### type GoProxy

```go
type GoProxy struct {
}
```

GoProxy is an API Gateway/Reverse Proxy and http.ServeMux/http.Handler

#### func  NewGoProxy

```go
func NewGoProxy(config *config.Config) *GoProxy
```
NewGoProxy registers a new reverseproxy handler for each provided config with
the specified path prefix

#### func (*GoProxy) ServeHTTP

```go
func (g *GoProxy) ServeHTTP(w http.ResponseWriter, r *http.Request)
```
