package main

import (
	"github.com/autom8ter/api/go/api"
	"github.com/autom8ter/goproxy"
	"github.com/autom8ter/goproxy/config"
	"net/http"
	"os"
)

var BaseURL = "https://api.stripe.com/v1/customers"

var proxy = goproxy.NewGoProxy(&config.Config{
	TargetUrl:  BaseURL,
	WebHookURL: os.Getenv("WEBHOOK"),
})

func main() {
	if err := http.ListenAndServe(":8080", proxy); err != nil {
		api.Util.Entry().Fatalln(err.Error())
	}
}
