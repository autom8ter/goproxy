package middleware

import (
	"github.com/gorilla/mux"
	"net/http"
)

//ResponseWare is a function used to modify the response of a reverse proxy
type ResponseWare func(func(response *http.Response) error) func(response *http.Response) error

//RequestWare is a function used to modify the incoming request of a reverse proxy from a client
type RequestWare func(func(req *http.Request)) func(req *http.Request)

//TransportWare is a function used to modify the http RoundTripper that is used by a reverse proxy. The default RoundTripper is initially http.DefaultTransport
type TransportWare func(tripper http.RoundTripper) http.RoundTripper

//RouterWare is a function used to modify the mux
type RouterWare func(r *mux.Router) *mux.Router
