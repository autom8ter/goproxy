package goproxy_test

import (
	"github.com/autom8ter/goproxy"
	"log"
	"net/http"
	"os"
	"testing"
)

var twilioTarget = "https://api.twilio.com/2010-04-01/Accounts/" + os.Getenv("TWILIO_ACCOUNT_SID")

func TestNewGoProxy(t *testing.T) {
	Config := goproxy.ProxyConfig{
		"/Messages.json": &goproxy.Config{
			TargetUrl: twilioTarget,
			Username:  os.Getenv("TWILIO_API_KEY"),
			Password:  os.Getenv("TWILIO_API_SECRET"),
			Headers: map[string]string{
				"Header_Key": "Header Value",
			},
		},
		"/Calls.json": &goproxy.Config{
			TargetUrl: twilioTarget,
			Username:  os.Getenv("TWILIO_API_KEY"),
			Password:  os.Getenv("TWILIO_API_SECRET"),
			Headers: map[string]string{
				"Header_Key": "Header Value",
			},
			FormValues: map[string]string{
				"Form_Key": "Form Value",
			},
		},
	}

	gProxy := goproxy.New(Config)
	if gProxy == nil {
		t.Fatal("registered nil goproxy")
	}
	if gProxy.Router == nil {
		t.Fatal("registered nil serveMux")
	}
	for k, v := range gProxy.Proxies() {
		if v == nil {
			t.Fatalf("registered nil reverse proxy with prefix: %s", k)
		}
		if v.Director == nil {
			t.Fatal("registered nil reverse proxy director")
		}
	}
	if err := http.ListenAndServe(":8080", gProxy); err != nil {
		log.Fatalln(err.Error())
	}
}
