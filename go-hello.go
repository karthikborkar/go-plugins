/*
A "hello world" plugin in Go,
which reads a request header and sets a response header.
*/
package main

import (
	"strings"

	"github.com/Kong/go-pdk"
)

// Config ...
type Config struct {
	Message string
}

// New ...
func New() interface{} {
	return &Config{}
}

var (
	defaultUpstream = "europe_cluster"
)

// GetHeaders gets the Value for Header X-Country
func GetHeaders(kong *pdk.PDK) (string, string, error) {

	country, err := kong.Request.GetHeader("X-Country")
	if err != nil {
		country = ""
		kong.Log.Err(err.Error())
	}

	region, err := kong.Request.GetHeader("X-Regione")
	if err != nil {
		region = ""
		kong.Log.Err(err.Error())
	}
	return strings.ToLower(country), strings.ToLower(region), err
}

// GetUpstream gets the dynamic upstream for Header X-Country
func GetUpstream(kong *pdk.PDK) string {

	// mapping of country, region, and upstream
	// can be kept in DB and we can fetch it from
	// there was well. it can be more dynamic if needed
	country, region, err := GetHeaders(kong)
	if err != nil || (country == "" && region == "") {
		return defaultUpstream
	} else if country == "italy" && region == "abruzzo" {
		return "italy_cluster"
	} else {
		return defaultUpstream
	}
}

// Access https://docs.konghq.com/1.0.x/plugin-development/custom-logic/
func (conf Config) Access(kong *pdk.PDK) {

	upstream := GetUpstream(kong)

	// https://docs.konghq.com/1.2.x/pdk/kong.service/#kongserviceset_upstreamhost
	err := kong.Service.SetUpstream(upstream)
	if err != nil {
		kong.Log.Err(err.Error())
	}

	kong.Response.SetHeader("x-kong-upstream", upstream)
}
