package main

import (
	"net/http"

	"github.com/cfrye2000/productPromisedEventMS/external/github.com/gorilla/mux"
	"github.com/cfrye2000/productPromisedEventMS/external/github.com/rs/cors"
)

//NewRouter  Mutex router for the application
func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		//first the actual application handler
		handler = route.HandlerFunc

		//now wrap that with a logger
		if !route.SkipLog {
			handler = Logger(handler, route.Name)
		}

		//add throttler handler
		if route.Throttle {
			handler = Throttler(handler)
		}

		//if this handler requires authentication first, wrap everything with an Authenticator
		if route.Authenticate {
			handler = Authenticater(handler, route.SecurityLevel)
		}

		//if ssl is required, test for that with a finall wrapper
		handler = SslTester(handler)

		//finally add a CORS handler for cross domain access
		c := cors.New(cors.Options{
			AllowedOrigins:     []string{"*"},
			AllowCredentials:   true,
			AllowedMethods:     []string{"GET", "POST", "DELETE", "PATCH", "PUT"},
			AllowedHeaders:     []string{"Authorization", "Accept", "Content-Type", "User-Credential", "X-Api-Key"},
			OptionsPassthrough: false,
			Debug:              false,
		})
		handler = c.Handler(handler)

		router.
			Methods(route.Method, "OPTIONS").
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}

	return router
}
