# goproxy
--
    import "github.com/autom8ter/goproxy"

## Example

```go

var BaseURL = "https://api.stripe.com/v1/customers"

var proxy = goproxy.NewGoProxy(&config.Config{
	TargetUrl:           BaseURL,
	Secret:              os.Getenv("SECRET"),//used for signing json web tokens
})

func main() {
	proxy.ListenAndServe(":8080")
}

```
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

#### func (*GoProxy) ListenAndServe

```go
func (g *GoProxy) ListenAndServe(addr string)
```

#### func (*GoProxy) ServeHTTP

```go
func (g *GoProxy) ServeHTTP(w http.ResponseWriter, r *http.Request)
```
