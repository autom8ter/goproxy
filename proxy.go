package goproxy

import (
	"github.com/autom8ter/goproxy/config"
	"github.com/autom8ter/goproxy/httputil"
	"github.com/autom8ter/goproxy/util"
	"net/http"
)

//GoProxy is a configurable single-target reverse-proxy HTTP handler compatible with the net/http http.Handler interface
type GoProxy struct {
	r *httputil.ReverseProxy
}

func (g *GoProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.r.ServeHTTP(w, r)
}

//NewGoProxy registers a new reverseproxy handler for each provided config with the specified path prefix
func NewGoProxy(config *config.Config) *GoProxy {
	if err := util.Handle.Validate(config); err != nil {
		util.Handle.Entry().Fatalln(err.Error())
	}
	return &GoProxy{
		r: &httputil.ReverseProxy{
			Director:         config.DirectorFunc(),
			Transport:        http.DefaultTransport,
			FlushInterval:    config.FlushInterval,
			ErrorLog:         config.Entry(),
			ResponseCallback: config.ResponseCallback(),
			ErrorHandler:     nil,
		},
	}
}
