package main

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"

	eventsMod "github.com/cfrye2000/productPromisedEventMS/events"
)

func init() {

	route := Route{
		"SubmitEvent",
		"POST",
		"/productpromisedevents",
		SubmitEvent,
		0,
		true,
		true,
		false,
	}

	routes = append(routes, route)

}

//SubmitEvent  Submit a Product Ordered Event
func SubmitEvent(w http.ResponseWriter, r *http.Request) {
	//get the post body

	p, readError := ioutil.ReadAll(r.Body)

	if readError != nil && readError != io.EOF {
		WriteResponse(w, r, http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Text: "Can not read posted image information"})
		return
	}

	contentType := getContentFormat(r)

	var events []eventsMod.ProductPromisedEvent
	var xmlRoot eventsMod.RootLevel
	var err error

	if contentType == "text/json" {
		err = json.Unmarshal(p, &events)
	} else {
		err = xml.Unmarshal(p, &xmlRoot)
		events = xmlRoot.ProductPromisedEvents
	}

	if err != nil {
		WriteResponse(w, r, http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Text: "Can not marshal posted event information"})
		return
	}

	if valid, m := validateEvents(events); !valid {
		WriteResponse(w, r, http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Text: "Need to specify all the required information for an event: " + m})
	} else {
		if response, err := eventsMod.SubmitEvent(events, psClient, ctx); err != nil {
			WriteResponse(w, r, http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Text: err.Error()})
		} else {
			WriteResponse(w, r, http.StatusOK, response)
		}
	}
}

func validateEvents(events []eventsMod.ProductPromisedEvent) (bool, string) {

	if len(events) > 5 {
		return false, "Too many events.  5 is the maximum"
	}

	for _, val := range validators {
		for _, e := range events {
			if v, err := val.ValidateFunc(e); !v {
				return v, err
			}
		}
	}

	return true, ""

}

func getContentFormat(r *http.Request) string {
	txt := r.Header.Get("Content-Type")
	var resp string
	if txt == "text/xml" {
		resp = "text/xml"
	} else {
		resp = "text/json"
	}

	return resp
}
