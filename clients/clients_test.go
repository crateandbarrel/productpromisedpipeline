package clients

import (
	"testing"
)

func TestClientRetrieval(t *testing.T) {
	var client Client
	if c, ok := Clients["chrislong"]; !ok {
		t.Errorf("Error retrieving client")
	} else {
		client = c
	}

	if client.Secret != "doggie" {
		t.Errorf("Error accessing client")
	}
}
