/*
A "hello world" plugin in Go,
which reads a request header and sets a response header.
*/
package main

import (
	"fmt"

	"github.com/Kong/go-pdk"
)

type Config struct {
	Message string
}

func New() interface{} {
	return &Config{}
}

// getCountryHeader gets the Value for Header X-Country
func getCountryHeader(kong *pdk.PDK) (string, error) {
	country, err := kong.Request.GetHeader("X-Country")
	if err != nil {
		country = ""
		kong.Log.Err(err.Error())
	}
	return country, err
}

// getUpstream gets the dynamic upstream for Header X-Country
func getUpstream(kong *pdk.PDK) string {
	defaultUpstream := "europe_cluster"

	country, err := getCountryHeader(kong)
	if err != nil || country == "" {
		return defaultUpstream
	}

	if country == "italy" {
		return "italy_cluster"
	}

	return defaultUpstream
}

func (conf Config) Access(kong *pdk.PDK) {

	upstream := getUpstream(kong)
	//kong.ServiceRequest.SetHeader("host", upstream)

	// https://docs.konghq.com/1.2.x/pdk/kong.service/#kongserviceset_upstreamhost
	//kong.Service.SetUpstream(upstream)

	kong.Response.SetHeader("x-hello-from-go", fmt.Sprintf("Go says %s to %s", "message", upstream))
}
