package goproxy

import (
	"github.com/auth0/go-jwt-middleware"
	"github.com/autom8ter/goproxy/config"
	"github.com/autom8ter/goproxy/httputil"
	"github.com/autom8ter/goproxy/util"
	"github.com/dgrijalva/jwt-go"
	"net/http"
)

//GoProxy is a configurable single-target reverse-proxy HTTP handler compatible with the net/http http.Handler interface
type GoProxy struct {
	r      *httputil.ReverseProxy
	auth   *jwtmiddleware.JWTMiddleware
	config *config.Config
}

func (g *GoProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := g.auth.CheckJWT(w, r)
	if err != nil {

		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	g.r.ServeHTTP(w, r)
}

//NewGoProxy registers a new reverseproxy handler for each provided config with the specified path prefix
func NewGoProxy(config *config.Config) *GoProxy {
	if err := util.Handle.Validate(config); err != nil {
		util.Handle.Entry().Fatalln(err.Error())
	}
	j := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Client_Secret), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	return &GoProxy{
		auth:   j,
		config: config,
		r: &httputil.ReverseProxy{
			Director:         config.DirectorFunc(),
			Transport:        http.DefaultTransport,
			FlushInterval:    config.FlushInterval,
			ErrorLog:         config.Entry(),
			ResponseCallback: config.ResponseCallback(),
		},
	}
}
