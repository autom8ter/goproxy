package config_test

import (
	"fmt"
	"log"
	"net/url"
	"testing"
)

func Test(t *testing.T) {
	u, err := url.Parse("https://api.twilio.com/2010-04-01/Calls.json")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Full: %v\n", u)
	fmt.Printf("User: %v\n", u.User)
	fmt.Printf("Host: %v\n", u.Host)
	fmt.Printf("ForcQuery: %v\n", u.ForceQuery)
	fmt.Printf("Scheme: %v\n", u.Scheme)
	fmt.Printf("RawPath: %v\n", u.RawPath)
	fmt.Printf("RawQuery: %v\n", u.RawQuery)
	fmt.Printf("Path: %v\n", u.Path)
	fmt.Printf("Fragment: %v\n", u.Fragment)
	fmt.Printf("Opaque: %v\n", u.Opaque)

}
