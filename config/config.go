package config

import (
	"github.com/autom8ter/goproxy/util"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"time"
)

//Config is used to configure a reverse proxy handler(one route)
type Config struct {
	TargetUrl           string `validate:"required"`
	Secret              string `validate:"required"`
	Headers             map[string]string
	FormValues          map[string]string
	FlushInterval       time.Duration
	ResponseCallbackURL string
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

func (c *Config) ResponseCallback() func(r *http.Response) error {
	if c.ResponseCallbackURL == "" {
		return nil
	}
	return func(r *http.Response) error {
		u, err := url.Parse(c.ResponseCallbackURL)
		if err != nil {
			util.Handle.Entry().Fatalln("failed to parse response callback url", err.Error())
		}

		resp, err := http.DefaultClient.Do(&http.Request{
			Method:     "POST",
			URL:        u,
			Header:     cloneHeader(r.Header),
			Body:       r.Body,
			Trailer:    nil,
			RemoteAddr: r.Request.RemoteAddr,
		})
		r = resp
		if err != nil {
			return err
		}
		return nil
	}
}

func (c *Config) Entry() *logrus.Entry {
	return util.Handle.Entry()
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

func cloneHeader(h http.Header) http.Header {
	h2 := make(http.Header, len(h))
	for k, vv := range h {
		vv2 := make([]string, len(vv))
		copy(vv2, vv)
		h2[k] = vv2
	}
	return h2
}
