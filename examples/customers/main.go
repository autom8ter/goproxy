package main

import (
	"github.com/autom8ter/goproxy"
	"github.com/autom8ter/goproxy/config"
	"github.com/autom8ter/goproxy/util"
	"net/http"
	"os"
)

var BaseURL = "https://api.stripe.com/v1/customers"

var proxy = goproxy.NewGoProxy(&config.Config{
	TargetUrl:           BaseURL,
	Secret:              os.Getenv("SECRET"),
	ResponseCallbackURL: os.Getenv("CALLBACK"),
})

func main() {
	util.Handle.Entry().Println("starting proxy on :8081")
	if err := http.ListenAndServe(":8081", proxy); err != nil {
		util.Handle.Entry().Fatalln(err.Error())
	}
}
