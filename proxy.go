package goproxy

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/http/pprof"
	"net/url"
	"strings"
)

//GoProxy is an API Gateway/Reverse Proxy and http.ServeMux/http.Handler
type GoProxy struct {
	*http.ServeMux
	proxies map[string]*httputil.ReverseProxy
}

//ResponseMiddleware is a function used to modify the response of a reverse proxy
type ResponseMiddleware func(func(response *http.Response) error) func(response *http.Response) error

//RequestMiddleware is a function used to modify the incoming request of a reverse proxy from a client
type RequestMiddleware func(func(req *http.Request)) func(req *http.Request)

//TransportMiddleware is a function used to modify the http RoundTripper that is used by a reverse proxy. The default RoundTripper is initially http.DefaultTransport
type TransportMiddleware func(tripper http.RoundTripper) http.RoundTripper

//ProxyConfig is used to configure GoProxies reverse proxies
type ProxyConfig struct {
	PathPrefix string
	TargetUrl  string
	Username   string
	Password   string
	Headers    map[string]string
}

//NewGoProxy registers a new reverseproxy for each provided ProxyConfig
func NewGoProxy(configs ...*ProxyConfig) *GoProxy {
	g := &GoProxy{
		ServeMux: http.NewServeMux(),
		proxies:  make(map[string]*httputil.ReverseProxy),
	}
	for _, c := range configs {
		g.proxies[c.PathPrefix] = &httputil.ReverseProxy{
			Director: g.directorFunc(c),
		}
	}
	for path, prox := range g.proxies {
		g.Handle(path, prox)
	}
	return g
}

func (g *GoProxy) directorFunc(config *ProxyConfig) func(req *http.Request) {
	target, err := url.Parse(config.TargetUrl)
	if err != nil {
		log.Fatalln(err.Error())
	}
	targetQuery := target.RawQuery
	return func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		if config.Username != "" && config.Password != "" {
			req.SetBasicAuth(config.Username, config.Password)
		}
		if config.Headers != nil {
			for k, v := range config.Headers {
				req.Header.Set(k, v)
			}
		}
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

//ModifyResponses takes a Response Middleware function, traverses each registered reverse proxy, and modifies the http response it sends to the client
func (g *GoProxy) ModifyResponses(middleware ResponseMiddleware) {
	for _, prox := range g.proxies {
		prox.ModifyResponse = middleware(prox.ModifyResponse)
	}
}

//ModifyResponses takes a Request Middleware function, traverses each registered reverse proxy, and modifies the http request it sends to its target prior to sending
func (g *GoProxy) ModifyRequests(middleware RequestMiddleware) {
	for _, prox := range g.proxies {
		prox.Director = middleware(prox.Director)
	}
}

//ModifyResponses takes a Transport Middleware function, traverses each registered reverse proxy, and modifies the http roundtripper it uses
func (g *GoProxy) ModifyTransport(middleware TransportMiddleware) {
	for _, prox := range g.proxies {
		prox.Transport = middleware(prox.Transport)
	}
}

//Proxies returns all registered reverse proxies as a map of prefix:reverse proxy
func (g *GoProxy) Proxies() map[string]*httputil.ReverseProxy {
	return g.proxies
}

//GetProxy returns the reverse proxy with the registered prefix
func (g *GoProxy) GetProxy(prefix string) *httputil.ReverseProxy {
	return g.proxies[prefix]
}

//AsHandlerFunc converts a GoProxy to an http.HandlerFunc for convenience
func (g *GoProxy) AsHandlerFunc() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		g.ServeHTTP(writer, request)
	}
}

func (g *GoProxy) WithPprof() *GoProxy {
	g.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	g.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	g.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	g.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	g.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	return g
}

//ListenAndServe starts an http server on the given address with all of the registered reverse proxies
func (g *GoProxy) ListenAndServe(addr string) error {
	for prefix, _ := range g.proxies {
		fmt.Println("<--------------------------------PROXIED-------------------------------->")
		fmt.Printf("Prefix-------------> %s\n", prefix)
	}
	log.Printf("Starting GoProxy Server, Address: %s", addr)
	return http.ListenAndServe(addr, g)
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
