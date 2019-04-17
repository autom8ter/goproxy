package twilio

import (
	"github.com/autom8ter/goproxy"
	"net/http"
	"os"
)

var acc = os.Getenv("TWILIO_ACCOUNT")
var BaseURL = "https://api.twilio.com/2010-04-01" + "/Accounts/" + acc

func TwilioHandler(w http.ResponseWriter, r *http.Request) {
	goproxy.New(&goproxy.Config{
		PathPrefix: "/call",
		TargetUrl:  BaseURL + "/Calls.json",
		Username:   acc,
		Password:   os.Getenv("TWILIO_KEY"),
	},
		&goproxy.Config{
			PathPrefix: "/sms",
			TargetUrl:  BaseURL + "/Messages.json",
		},
	).ServeHTTP(w, r)
}
