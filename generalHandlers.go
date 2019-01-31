package main

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/cfrye2000/productPromisedEventMS/clients"
	"github.com/cfrye2000/productPromisedEventMS/external/github.com/gorilla/context"
	"github.com/cfrye2000/productPromisedEventMS/oauth"
	"github.com/cfrye2000/productPromisedEventMS/oauth2"
)

func init() {

	route := Route{
		"GetToken",
		"POST",
		"/oauth2/token",
		GetToken,
		0,
		false,
		false,
		false,
	}

	routes = append(routes, route)

}

//GetToken  Get an Oauth2 token
func GetToken(w http.ResponseWriter, r *http.Request) {
	client := clients.Client{}
	var err error

	txt := r.Header.Get("Authorization")

	if !strings.Contains(txt, "Basic") {
		e := errors.New("Unauthorized Access")
		WriteResponse(w, r, http.StatusUnauthorized, errorResponse{Code: http.StatusUnauthorized, Text: e.Error()})
		return
	}

	if client, err = oauth.AuthenticateClient(txt); err != nil {
		WriteResponse(w, r, http.StatusUnauthorized, errorResponse{Code: http.StatusUnauthorized, Text: err.Error()})
		return
	}

	context.Set(r, ClientID, client.ClientID)

	//now get the post body

	p := make([]byte, r.ContentLength)
	if _, err := r.Body.Read(p); err != nil && err != io.EOF {
		WriteResponse(w, r, http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Text: "Can not read Authentication Body"})
		return
	}

	bodyString := string(p[:])
	//first separate grant_type from refresh_token
	bodyArray := strings.Split(bodyString, "\n")

	if !strings.Contains(bodyArray[0], "grant_type") {
		WriteResponse(w, r, http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Text: "Invalid Authentication Body"})
		return
	}

	grantString := strings.TrimPrefix(bodyArray[0], "grant_type=")
	grantString = strings.TrimSpace(grantString)

	var response oauth2.TokenResponse
	err = nil

	if grantString == "refresh_token" {
		//test to make sure they supply a token
		if !strings.Contains(bodyArray[1], "refresh_token=") {
			WriteResponse(w, r, http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Text: "No token to refresh given"})
			return
		}
		tokenString := strings.TrimPrefix(bodyArray[1], "refresh_token=")
		tokenString = strings.TrimSpace(tokenString)
		if response, err = oauth2.RefreshToken(client.ClientID, tokenString, grantString, client.SecurityLevel, client.HitsPerMinute, client.Expiry); err != nil {
			WriteResponse(w, r, http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Text: err.Error()})
			return
		}
	} else if grantString == "client_credentials" {
		if response, err = oauth2.RequestToken(client.ClientID, grantString, client.SecurityLevel, client.HitsPerMinute, client.Expiry); err != nil {
			WriteResponse(w, r, http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Text: "Can not issue tokens"})
			return
		}
	} else {
		WriteResponse(w, r, http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Text: "Invalid Authentication Body"})
		return
	}

	WriteResponse(w, r, http.StatusOK, response)
}
