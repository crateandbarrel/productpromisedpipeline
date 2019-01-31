package main

import (
	"encoding/json"
	"log"
	"net/http"
)

//SslTester  Router that wraps other routers for Oauth2 authenticating
func SslTester(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if secure && r.TLS == nil {
			log.Printf(
				"%d\t%s\t%s\t%s\t%s",
				http.StatusBadRequest,
				r.Method,
				r.RequestURI,
				"ssl required",
				r.RemoteAddr,
			)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(errorResponse{Code: http.StatusUnauthorized, Text: "ssl required"}); err != nil {
				panic(err)
			}
		} else {
			inner.ServeHTTP(w, r)
		}
	})
}
