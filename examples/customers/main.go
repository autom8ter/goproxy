package main

import (
	"github.com/autom8ter/goproxy"
	"github.com/autom8ter/goproxy/config"
	"os"
)

var BaseURL = "https://api.stripe.com/v1/customers"

var proxy = goproxy.NewGoProxy(&config.Config{
	TargetUrl:  BaseURL,
	Secret:     os.Getenv("SECRET"),
	WebHookURL: os.Getenv("WEBHOOK"),
})

func main() {
	proxy.ListenAndServe(":8080")
}
