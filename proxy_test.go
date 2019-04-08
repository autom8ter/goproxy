package goproxy_test

import (
	"github.com/autom8ter/goproxy"
	"os"
	"testing"
)

func TestNewGoProxy(t *testing.T) {
	var (
		twilioTarget = "https://api.twilio.com/2010-04-01/Accounts/" + os.Getenv("TWILIO_ACCOUNT_SID")
		Configs      = []*goproxy.ProxyConfig{
			{
				PathPrefix: "/Messages.json",
				TargetUrl:  twilioTarget,
				Username:   os.Getenv("TWILIO_API_KEY"),
				Password:   os.Getenv("TWILIO_API_SECRET"),
				Headers: map[string]string{
					"Header_Key": "Header Value",
				},
			},
			{
				PathPrefix: "/Calls.json",
				TargetUrl:  twilioTarget,
				Username:   os.Getenv("TWILIO_API_KEY"),
				Password:   os.Getenv("TWILIO_API_SECRET"),
				Headers: map[string]string{
					"Header_Key": "Header Value",
				},
			},
		}
	)

	gProxy := goproxy.NewGoProxy(Configs...)
	if gProxy == nil {
		t.Fatal("registered nil goproxy")
	}
	if gProxy.ServeMux == nil {
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
}
