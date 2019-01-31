package main

import "net/http"

//Route structure of a route used by gorilla
type Route struct {
	Name          string
	Method        string
	Pattern       string
	HandlerFunc   http.HandlerFunc
	SecurityLevel int
	Authenticate  bool
	Throttle      bool
	SkipLog       bool
}

//Routes list of all the routes
type Routes []Route

var routes Routes
