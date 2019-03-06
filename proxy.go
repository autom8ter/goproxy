package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type GoProxy struct {
	*http.ServeMux
}

type ProxyConfig struct {
	PathPrefix string
	TargetUrl  string
	Username   string
	Password   string
}

func NewGoProxy(configs ...*ProxyConfig) *GoProxy {
	g := &GoProxy{
		ServeMux: http.NewServeMux(),
	}
	for _, c := range configs {
		g.Handle(c.PathPrefix, &httputil.ReverseProxy{
			Director:       g.directorFunc(c),
			ModifyResponse: g.responderFunc(),
		})
	}
	return g
}

func (g *GoProxy) directorFunc(config *ProxyConfig) func(req *http.Request) {
	target, err := url.Parse(config.TargetUrl)
	if err != nil {
		log.Fatalln()
	}
	targetQuery := target.RawQuery
	return func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		req.SetBasicAuth(config.Username, config.Password)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}
}

func (g *GoProxy) responderFunc() func(response *http.Response) error {
	return func(response *http.Response) error {
		response.Header = nil
		response.Request.Header = nil
		return nil
	}
}

func (g *GoProxy) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, g)
}
