package middleware

import (
	"fmt"
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	"net/http"
	"net/http/pprof"
)

type RequestLog struct {
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

//ResponseWare is a function used to modify the response of a reverse proxy
type ResponseWare func(func(response *http.Response) error) func(response *http.Response) error

//RequestWare is a function used to modify the incoming request of a reverse proxy from a client
type RequestWare func(func(req *http.Request)) func(req *http.Request)

//TransportWare is a function used to modify the http RoundTripper that is used by a reverse proxy. The default RoundTripper is initially http.DefaultTransport
type TransportWare func(tripper http.RoundTripper) http.RoundTripper

//RouterWare is a function used to modify the mux
type RouterWare func(r *mux.Router) *mux.Router

func WithJWT(secret string) RouterWare {
	return func(r *mux.Router) *mux.Router {
		j := jwtmiddleware.New(jwtmiddleware.Options{
			ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			},
			SigningMethod: jwt.SigningMethodHS256,
		})
		j.Handler(r)
		return r
	}
}

func WithMetrics() RouterWare {
	return func(r *mux.Router) *mux.Router {
		var (
			inFlightGauge = prometheus.NewGauge(prometheus.GaugeOpts{
				Name: "in_flight_requests",
				Help: "A gauge of requests currently being served by the wrapped handler.",
			})

			counter = prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Name: "api_requests_total",
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
					Buckets: []float64{.025, .05, .1, .25, .5, 1},
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
		_ = r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
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
		fmt.Println("registered handler: ", "/metrics")
		r.Handle("/metrics", promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}))
		return r
	}
}

func WithProf() RouterWare {
	return func(r *mux.Router) *mux.Router {
		r.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
		r.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
		r.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
		r.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
		r.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
		return r
	}
}

func WithCORS(options cors.Options) RouterWare {
	return func(r *mux.Router) *mux.Router {
		cors.New(options).Handler(r)
		return r
	}
}
