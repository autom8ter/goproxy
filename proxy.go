package goproxy

import (
	"github.com/autom8ter/goproxy/middleware"
	"github.com/autom8ter/objectify"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

var util = objectify.Default()

//GoProxy is an API Gateway/Reverse Proxy and http.ServeMux/http.Handler
type GoProxy struct {
	*mux.Router
	proxies map[string]*httputil.ReverseProxy
}

//Config is used to configure a reverse proxy handler(one route)
type Config struct {
	PathPrefix string `validate:"required"`
	TargetUrl  string `validate:"required"`
	Username   string
	Password   string
	Headers    map[string]string
	FormValues map[string]string
}

//New registers a new reverseproxy handler for each provided config with the specified path prefix
func New(configs ...*Config) *GoProxy {
	if len(configs) == 0 {
		util.Entry().Warnln("zero configs passed in creation of GoProxy")
	}
	proxy := &GoProxy{
		Router:  mux.NewRouter(),
		proxies: make(map[string]*httputil.ReverseProxy),
	}
	for _, v := range configs {
		if err := util.Validate(v); err != nil {
			util.Entry().Fatalln(err.Error())
		}
		proxy.proxies[v.PathPrefix] = &httputil.ReverseProxy{
			Director: proxy.directorFunc(v),
		}
	}
	for path, prox := range proxy.proxies {
		proxy.Handle(path, prox)
	}
	return proxy
}

//NewSecure registers a new secure reverseproxy for each provided configs. It is the same as New, except with CORS options and a
// JWT middleware that checks for a signed bearer token
func NewSecure(secret string, opts cors.Options, configs ...*Config) *GoProxy {
	if len(configs) == 0 {
		util.Entry().Warnln("zero configs passed in creation of GoProxy")
	}
	proxy := &GoProxy{
		Router:  mux.NewRouter(),
		proxies: make(map[string]*httputil.ReverseProxy),
	}
	for _, v := range configs {
		if err := util.Validate(v); err != nil {
			util.Entry().Fatalln(err.Error())
		}
		proxy.proxies[v.PathPrefix] = &httputil.ReverseProxy{
			Director: proxy.directorFunc(v),
		}
	}
	for path, prox := range proxy.proxies {
		proxy.Handle(path, prox)
	}
	proxy.Router = middleware.WithJWT(secret)(proxy.Router)
	proxy.Router = middleware.WithCORS(opts)(proxy.Router)
	return proxy
}

func (g *GoProxy) directorFunc(config *Config) func(req *http.Request) {
	target, err := url.Parse(config.TargetUrl)
	if err != nil {
		util.Entry().Fatalln(err.Error())
	}
	targetQuery := target.RawQuery
	return func(req *http.Request) {
		start := time.Now()
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = util.SingleJoiningSlash(target.Path, req.URL.Path)
		if config.Username != "" && config.Password != "" {
			req.SetBasicAuth(config.Username, config.Password)
		}
		if config.Headers != nil {
			for k, v := range config.Headers {
				req.Header.Set(k, v)
			}
		}
		if config.FormValues != nil {
			for k, v := range config.FormValues {
				req.Form.Set(k, v)
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

		util.Entry().Debugf("proxied request: %s\n", util.MarshalJSON(&middleware.RequestLog{
			Received:  util.HumanizeTime(start),
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

//ModifyResponses takes a Response Middleware function, traverses each registered reverse proxy, and modifies the http response it sends to the client
func (g *GoProxy) ModifyResponses(middleware middleware.ResponseWare) {
	for _, prox := range g.proxies {
		prox.ModifyResponse = middleware(prox.ModifyResponse)
	}
}

//ModifyResponses takes a Request Middleware function, traverses each registered reverse proxy, and modifies the http request it sends to its target prior to sending
func (g *GoProxy) ModifyRequests(middleware middleware.RequestWare) {
	for _, prox := range g.proxies {
		prox.Director = middleware(prox.Director)
	}
}

//ModifyResponses takes a Transport Middleware function, traverses each registered reverse proxy, and modifies the http roundtripper it uses
func (g *GoProxy) ModifyTransport(middleware middleware.TransportWare) {
	for _, prox := range g.proxies {
		prox.Transport = middleware(prox.Transport)
	}
}

//ModifyRouter takes a router middleware function and wraps the proxies router
func (g *GoProxy) ModifyRouter(middleware middleware.RouterWare) {
	g.Router = middleware(g.Router)
}

//WalkPaths walks registered mux paths and modifies them
func (g *GoProxy) WalkPaths(fns ...mux.WalkFunc) error {
	for _, v := range fns {
		if err := g.Router.Walk(v); err != nil {
			return err
		}
	}
	return nil
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

//ListenAndServe starts the GoProxy server on the specified address
func (g *GoProxy) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, g)
}
