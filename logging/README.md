# logging
--
    import "github.com/autom8ter/goproxy/logging"


## Usage

#### type Request

```go
type Request struct {
	Received time.Time `json:"received"`
	Method   string    `json:"method"`
	URL      string    `json:"url"`

	Body      string        `json:"body"`
	UserAgent string        `json:"user_agent"`
	Referer   string        `json:"referer"`
	Proto     string        `json:"proto"`
	RemoteIP  string        `json:"remote_ip"`
	Latency   time.Duration `json:"latency"`
}
```
