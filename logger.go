package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cfrye2000/productPromisedEventMS/external/github.com/gorilla/context"
)

func init() {
	log.SetOutput(os.Stdout)
}

//Logger  Router that wraps other routers for logging
func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		var clientID string

		if val, ok := context.GetOk(r, ClientID); ok {
			clientID = val.(string)
		} else {
			clientID = "unknown"
		}

		log.Printf(
			"\t%s\t%s\t%s\t%s\t%s\t%s\t%s",
			w.Header().Get("Status"),
			r.RemoteAddr,
			clientID,
			name,
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
		context.Clear(r)
	})
}
