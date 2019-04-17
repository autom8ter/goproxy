package goproxy

import (
	"github.com/autom8ter/goproxy/middleware"
	"github.com/autom8ter/objectify"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	"net/http"
	"net/http/httputil"
	"net/http/pprof"
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

//NewGoProxy registers a new reverseproxy handler for each provided config with the specified path prefix
func NewGoProxy(configs ...*Config) *GoProxy {
	if len(configs) == 0 {
		util.Entry().Warnln("zero configs passed in creation of GoProxy")
	}
	r := mux.NewRouter()
	proxies := make(map[string]*httputil.ReverseProxy)
	for _, v := range configs {
		if err := util.Validate(v); err != nil {
			util.Entry().Fatalln(err.Error())
		}
		proxies[v.PathPrefix] = &httputil.ReverseProxy{
			Director: directorFunc(v),
		}
	}
	for path, prox := range proxies {
		r.Handle(path, prox)
	}
	return &GoProxy{
		Router:  r,
		proxies: proxies,
	}
}

//NewSecureGoProxy registers a new secure reverseproxy for each provided configs. It is the same as New, except with CORS options and a
// JWT middleware that checks for a signed bearer token
func NewSecureGoProxy(secret string, opts cors.Options, configs ...*Config) *GoProxy {
	if len(configs) == 0 {
		util.Entry().Warnln("zero configs passed in creation of GoProxy")
	}
	r := mux.NewRouter()
	proxies := make(map[string]*httputil.ReverseProxy)
	for _, v := range configs {
		if err := util.Validate(v); err != nil {
			util.Entry().Fatalln(err.Error())
		}
		proxies[v.PathPrefix] = &httputil.ReverseProxy{
			Director: directorFunc(v),
		}
	}
	for path, prox := range proxies {
		r.Handle(path, prox)
	}
	return &GoProxy{
		Router:  r,
		proxies: proxies,
	}
}

func directorFunc(config *Config) func(req *http.Request) {
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
func (g *GoProxy) ResponseWare(middleware middleware.ResponseWare) {
	for _, prox := range g.proxies {
		prox.ModifyResponse = middleware(prox.ModifyResponse)
	}
}

//ModifyResponses takes a Request Middleware function, traverses each registered reverse proxy, and modifies the http request it sends to its target prior to sending
func (g *GoProxy) RequestWare(middleware middleware.RequestWare) {
	for _, prox := range g.proxies {
		prox.Director = middleware(prox.Director)
	}
}

//ModifyResponses takes a Transport Middleware function, traverses each registered reverse proxy, and modifies the http roundtripper it uses
func (g *GoProxy) TransportWare(middleware middleware.TransportWare) {
	for _, prox := range g.proxies {
		prox.Transport = middleware(prox.Transport)
	}
}

//Middleware wraps Goproxy with the provided middlewares
func (g *GoProxy) Middleware(middlewares ...mux.MiddlewareFunc) {
	g.Router.Use(middlewares...)
}

//WalkPaths walks registered mux paths
func (g *GoProxy) WalkPaths(walkfuncs ...mux.WalkFunc) error {
	for _, v := range walkfuncs {
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

//Registers prometheus metrics for: in_flight_requests, requests_total, request_duration_seconds, response_size_bytes,
func (g *GoProxy) WithMetrics() {
	var (
		inFlightGauge = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "in_flight_requests",
			Help: "A gauge of requests currently being served by the wrapped handler.",
		})

		counter = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "requests_total",
				Help: "A counter for requests to the wrapped handler.",
			},
			[]string{"code", "method"},
		)

		// duration is partitioned by the HTTP method and handler. It uses custom
		// buckets based on the expected request duration.
		duration = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "request_duration_seconds",
				Help:    "A histogram of latencies for requests.",
				Buckets: []float64{.005, .005, .1, .25, .5, 1},
			},
			[]string{"handler", "method"},
		)

		// responseSize has no labels, making it a zero-dimensional
		// ObserverVec.
		responseSize = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "response_size_bytes",
				Help:    "A histogram of response sizes for requests.",
				Buckets: []float64{200, 500, 900, 1500},
			},
			[]string{},
		)
	)

	// Register all of the metrics in the standard registry.
	prometheus.MustRegister(inFlightGauge, counter, duration, responseSize)
	var chain http.Handler
	_ = g.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pth, _ := route.GetPathTemplate()
		chain = promhttp.InstrumentHandlerInFlight(inFlightGauge,
			promhttp.InstrumentHandlerDuration(duration.MustCurryWith(prometheus.Labels{"handler": pth}),
				promhttp.InstrumentHandlerCounter(counter,
					promhttp.InstrumentHandlerResponseSize(responseSize, route.GetHandler()),
				),
			),
		)
		route = route.Handler(chain)
		return nil
	})
	g.Handle("/metrics", promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}))
}

//registers all pprof handlers: /debug/pprof/, /debug/pprof/cmdline, /debug/pprof/profile, /debug/pprof/symbol, /debug/pprof/trace
func (g *GoProxy) WithPprof() {
	g.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	g.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	g.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	g.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	g.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
}
