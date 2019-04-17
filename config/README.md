# config
--
    import "github.com/autom8ter/goproxy/config"


## Usage

#### type Config

```go
type Config struct {
	TargetUrl           string `validate:"required"`
	Secret              string `validate:"required"`
	Headers             map[string]string
	FormValues          map[string]string
	FlushInterval       time.Duration
	ResponseCallbackURL string
}
```

Config is used to configure a reverse proxy handler(one route)

#### func (*Config) DirectorFunc

```go
func (c *Config) DirectorFunc() func(req *http.Request)
```

#### func (*Config) Entry

```go
func (c *Config) Entry() *logrus.Entry
```

#### func (*Config) JSONString

```go
func (c *Config) JSONString() string
```

#### func (*Config) ResponseCallback

```go
func (c *Config) ResponseCallback() func(r *http.Response) error
```
