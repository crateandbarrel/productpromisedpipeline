package main

import (
	"net/http"
)

func init() {

	route := Route{
		"Ping",
		"GET",
		"/ping",
		Ping,
		0,
		true,
		true,
		false,
	}

	routes = append(routes, route)

	route = Route{
		"HealthCheck",
		"GET",
		"/healthcheck",
		HealthCheck,
		0,
		false,
		false,
		true,
	}

	routes = append(routes, route)

	route = Route{
		"HealthCheck",
		"GET",
		"/productpromisedevents/healthcheck",
		HealthCheck,
		0,
		false,
		false,
		true,
	}

	routes = append(routes, route)

	route = Route{
		"HealthCheck",
		"GET",
		"/",
		HealthCheck,
		0,
		false,
		false,
		true,
	}

	routes = append(routes, route)

}

//Ping  Get Ping and Location based on IP address
func Ping(w http.ResponseWriter, r *http.Request) {

	ipAddress := r.RemoteAddr
	pong := Pong{}
	pong.IPAddress = ipAddress

	WriteResponse(w, r, http.StatusOK, pong)

}

//HealthCheck  Check the health of server
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	ipAddress := r.RemoteAddr
	healthcheck := HealthCheckResponse{}
	healthcheck.IPAddress = ipAddress
	WriteResponse(w, r, http.StatusOK, healthcheck)
}
