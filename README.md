# goproxy
--
    import "github.com/autom8ter/goproxy"


## Usage

#### type GoProxy

```go
type GoProxy struct {
}
```

GoProxy is a configurable single-target reverse-proxy HTTP handler compatible
with the net/http http.Handler interface

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
