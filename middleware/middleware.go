package middleware

import (
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/rs/cors"
	"net/http"
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
type HandlerWare func(h http.Handler) http.Handler

func WithJWT(secret string) HandlerWare {
	return func(h http.Handler) http.Handler {
		j := jwtmiddleware.New(jwtmiddleware.Options{
			ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			},
			SigningMethod: jwt.SigningMethodHS256,
		})
		j.Handler(h)
		return h
	}
}

func WithCORS(options cors.Options) HandlerWare {
	return func(h http.Handler) http.Handler {
		cors.New(options).Handler(h)
		return h
	}
}
