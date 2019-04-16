# GoProxy

[GoDoc](https://github.com/autom8ter/goproxy/blob/master/GODOC.md)

## Overview
GoProxy is a lightweight reverse proxy server written in Golang

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
	// An example target url to proxy to:
	twilioTarget = "https://api.twilio.com/2010-04-01/Accounts/" + os.Getenv("TWILIO_ACCOUNT_SID")
	
	// Register Config- this is path prefixes(for the router/mux) to GoProxy configuration
	Config     = goproxy.ProxyConfig{
	
    		"/Messages.json" : &goproxy.Config{
    			TargetUrl:  twilioTarget,
    			Username:   os.Getenv("TWILIO_API_KEY"),
    			Password:   os.Getenv("TWILIO_API_SECRET"),
    			Headers: map[string]string{
    				"Header_Key": "Header Value",
    			},
    			FormValues: map[string]string{
    				"Form_Key": "Form Value",
    			},
    		},
    		"/Calls.json" : &goproxy.Config{
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

)

func main() {
	var goProxy = goproxy.New(Config)
	if err := http.ListenAndServe(":8080", goProxy); err != nil {
    		log.Fatalln(err.Error())
    	}
}

```

### Client (curl)

```text
curl -X POST localhost:8080/Messages.json -d "To={SEND_TO_NUMBER" "From={SEND_FROM_NUMBER" "Body={SMS_BODY"
```

## Configuration

You may modify the incoming request with the following variables:
- target url
- username
- password
- headers
- form values

```text

//Config is used to configure GoProxies reverse proxies
type Config struct {
	TargetUrl  string `validate:"required"`
	Username   string
	Password   string
	Headers    map[string]string
	FormValues map[string]string
}
```