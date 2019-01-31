package main

import (
	"encoding/json"
	"net/http"

	"github.com/cfrye2000/productPromisedEventMS/external/github.com/gorilla/context"
	"github.com/cfrye2000/productPromisedEventMS/oauth"
)

type key int

//ClientID Used as a context Key in the request
const ClientID key = 0

//AccessToken Used as a context Key in the request
const AccessToken key = 1

//Authenticater  Router that wraps other routers for Oauth2 authenticating
func Authenticater(inner http.Handler, routeSecurityLevel int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if t, err := oauth.AuthenticateToken(w, r); err != nil {
			logErr{Code: http.StatusUnauthorized, RemoteAddr: r.RemoteAddr, ClientID: t.ClientID, Error: err.Error() + " (" + t.Token + ")", Method: r.Method, RequestURI: r.RequestURI}.writeErrorToLog()
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusUnauthorized)
			if err := json.NewEncoder(w).Encode(errorResponse{Code: http.StatusUnauthorized, Text: err.Error()}); err != nil {
				panic(err)
			}
		} else if t.SecurityLevel < routeSecurityLevel {
			logErr{Code: http.StatusUnauthorized, RemoteAddr: r.RemoteAddr, ClientID: t.ClientID, Error: "Not authorized to use this service", Method: r.Method, RequestURI: r.RequestURI}.writeErrorToLog()
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusUnauthorized)
			if err := json.NewEncoder(w).Encode(errorResponse{Code: http.StatusUnauthorized, Text: "Not authorized to use this service"}); err != nil {
				panic(err)
			}
		} else {
			context.Set(r, ClientID, t.ClientID)
			context.Set(r, AccessToken, t)
			inner.ServeHTTP(w, r)
		}
	})
}
