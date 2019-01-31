package oauth

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"

	"github.com/cfrye2000/productPromisedEventMS/clients"
	"github.com/cfrye2000/productPromisedEventMS/oauth2"
)

//AuthenticateToken Make sure there is a bearer token and that it is valid
func AuthenticateToken(w http.ResponseWriter, r *http.Request) (oauth2.AccessToken, error) {
	txt := r.Header.Get("Authorization")
	if !strings.Contains(txt, "Bearer") {
		return oauth2.AccessToken{}, errors.New("Misformed Header")
	}

	token := strings.TrimPrefix(txt, "Bearer")
	token = strings.TrimSpace(token)

	accessToken, err := oauth2.AuthToken(token)

	if err != nil {
		return accessToken, err
	}

	if _, err := isKnownClient(accessToken.ClientID); err != nil {
		return accessToken, err
	}

	return accessToken, nil
}

//AuthenticateClient used to see if this is a client we recognize
func AuthenticateClient(txt string) (clients.Client, error) {

	if !strings.Contains(txt, "Basic") {
		e := errors.New("Unauthorized Access")
		return clients.Client{}, e
	}

	authString := strings.TrimPrefix(txt, "Basic")
	authString = strings.TrimSpace(authString)

	var client clients.Client
	var returnError error

	data, err := base64.StdEncoding.DecodeString(authString)

	if err != nil {
		returnError = errors.New("Can not decode auth string")
	} else {
		//test if is a know client
		clientString := string(data[:])
		clientStringArray := strings.Split(clientString, ":")
		clientID := clientStringArray[0]
		c, err := isKnownClient(clientID)
		client = c
		if len(clientStringArray) < 2 || err != nil {
			returnError = errors.New("Unauthorized Access")
		} else {
			//test the password
			password := clientStringArray[1]
			if password != client.Secret {
				returnError = errors.New("Unauthorized Access")
			}
		}
	}

	return client, returnError
}

func isKnownClient(clientString string) (clients.Client, error) {
	var client clients.Client
	var returnError error
	//test if is a know client
	c, ok := clients.Clients[clientString]
	client = c
	if !ok || !c.Active {
		returnError = errors.New("Unauthorized Access")
	}
	return client, returnError
}
