package goproxy

import (
	"github.com/autom8ter/goproxy/logging"
	"github.com/autom8ter/goproxy/middleware"
	"github.com/autom8ter/objectify"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

var util = objectify.Default()

//GoProxy is an API Gateway/Reverse Proxy and http.ServeMux/http.Handler
type GoProxy struct {
	gmux *runtime.ServeMux
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
	IsGrpc     bool
}

//ProxyConfig configures the entire reverse proxy
type ProxyConfig struct {
	Configs []*Config
}

//New registers a new reverseproxy for each provided ProxyConfig
func New(config *ProxyConfig) *GoProxy {
	g := &GoProxy{
		Router:  mux.NewRouter(),
		proxies: make(map[string]*httputil.ReverseProxy),
	}
	for _, v := range config.Configs {
		if err := util.Validate(v); err != nil {
			util.Entry().Fatalln(err.Error())
		}
		g.proxies[v.PathPrefix] = &httputil.ReverseProxy{
			Director: g.directorFunc(v),
		}
	}
	for path, prox := range g.proxies {
		g.Handle(path, prox)
	}
	return g
}

func (g *GoProxy) directorFunc(config *Config) func(req *http.Request) {
	target, err := url.Parse(config.TargetUrl)
	if err != nil {
		log.Fatalln(err.Error())
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

		util.Entry().Debugf("proxied request: %s\n", util.MarshalJSON(&logging.Request{
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

func (g *GoProxy) handleGRPC(grpcServer *grpc.Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			g.ServeHTTP(w, r)
		}
	})
}

func (g *GoProxy) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, g)
}
