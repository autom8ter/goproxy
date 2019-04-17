# middleware
--
    import "github.com/autom8ter/goproxy/middleware"


## Usage

#### type RequestLog

```go
type RequestLog struct {
	Received  string `json:"received"`
	Method    string `json:"method"`
	URL       string `json:"url"`
	Body      string `json:"body"`
	UserAgent string `json:"user_agent"`
	Referer   string `json:"referer"`
	Proto     string `json:"proto"`
	RemoteIP  string `json:"remote_ip"`
	Latency   string `json:"latency"`
}
```


#### type RequestWare

```go
type RequestWare func(func(req *http.Request)) func(req *http.Request)
```

RequestWare is a function used to modify the incoming request of a reverse proxy
from a client

#### type ResponseWare

```go
type ResponseWare func(func(response *http.Response) error) func(response *http.Response) error
```

ResponseWare is a function used to modify the response of a reverse proxy

#### type RouterWare

```go
type RouterWare func(r *mux.Router) *mux.Router
```

RouterWare is a function used to modify the mux

#### func  WithJWT

```go
func WithJWT(secret string) RouterWare
```

#### func  WithMetrics

```go
func WithMetrics() RouterWare
```

#### func  WithProf

```go
func WithProf() RouterWare
```

#### type TransportWare

```go
type TransportWare func(tripper http.RoundTripper) http.RoundTripper
```

TransportWare is a function used to modify the http RoundTripper that is used by
a reverse proxy. The default RoundTripper is initially http.DefaultTransport
