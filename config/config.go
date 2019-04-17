package config

import (
	"github.com/autom8ter/goproxy/util"
	"net/http"
	"net/url"
	"time"
)

//Config is used to configure a reverse proxy handler(one route)
type Config struct {
	TargetUrl     string `validate:"required"`
	Username      string
	Password      string
	Headers       map[string]string
	FormValues    map[string]string
	FlushInterval time.Duration
}

func (c *Config) DirectorFunc() func(req *http.Request) {
	target, err := url.Parse(c.TargetUrl)
	if err != nil {
		util.Handle.Entry().Fatalln(err.Error())
	}
	targetQuery := target.RawQuery
	return func(req *http.Request) {
		start := time.Now()
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = target.Path
		if c.Username != "" || c.Password != "" {
			req.SetBasicAuth(c.Username, c.Password)
		}
		if c.Headers != nil {
			for k, v := range c.Headers {
				req.Header.Set(k, v)
			}
		}
		if c.FormValues != nil {
			for k, v := range c.FormValues {
				req.Form.Set(k, v)
			}
		}
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}

		util.Handle.Entry().Debugf("proxied request: %s\n", util.Handle.MarshalJSON(&requestLog{
			Received:  util.Handle.HumanizeTime(start),
			Method:    req.Method,
			URL:       req.URL.String(),
			UserAgent: req.UserAgent(),
			Referer:   req.Referer(),
			Proto:     req.Proto,
			RemoteIP:  req.RemoteAddr,
			Latency:   time.Since(start).String(),
		}))
	}
}

func (c *Config) JSONString() string {
	return string(util.Handle.MarshalJSON(c))
}

func (c *Config) PublishResponse() string {
	return string(util.Handle.MarshalJSON(c))
}

type requestLog struct {
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
