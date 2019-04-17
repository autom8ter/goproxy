# config
--
    import "github.com/autom8ter/goproxy/config"


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
