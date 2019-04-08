# GoProxy

## Overview
GoProxy is a lightweight(zero third party libraries) reverse proxy server written in Golang

It registers target urls and appends basic authentication to the inbound request so that you may authorize access to 
many different API endpoints from a single gateway.


## Example Usage (Twilio)

### Reverse Proxy
```go
package main

import (
	"github.com/autom8ter/goproxy"
	"log"
	"os"
)
var (
	twilioTarget = "https://api.twilio.com/2010-04-01/Accounts/" + os.Getenv("TWILIO_ACCOUNT_SID")
	Configs      = []*goproxy.ProxyConfig{
		{
        				PathPrefix: "/Messages.json",
        				TargetUrl:  twilioTarget,
        				Username:   os.Getenv("TWILIO_API_KEY"),
        				Password:   os.Getenv("TWILIO_API_SECRET"),
        				Headers: map[string]string{
        					"Header_Key" : "Header Value",
        				},
        			},
        			{
        				PathPrefix: "/Calls.json",
        				TargetUrl:  twilioTarget,
        				Username:   os.Getenv("TWILIO_API_KEY"),
        				Password:   os.Getenv("TWILIO_API_SECRET"),
        				Headers: map[string]string{
        					"Header_Key" : "Header Value",
        				},
        			},
	}
)

func main() {
	var goProxy = goproxy.NewGoProxy(Configs...)
	log.Fatalln(goProxy.ListenAndServe(":8080"))
}

```

### Client (curl)

```text
curl -X POST localhost:8080/Messages.json -d "To={SEND_TO_NUMBER" "From={SEND_FROM_NUMBER" "Body={SMS_BODY"
```