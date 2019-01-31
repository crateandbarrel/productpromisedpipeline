package main

import (
	"log"
)

type errorResponse struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

type successResponse struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

type countResponse struct {
	Count int    `json:"count"`
	Text  string `json:"text"`
}

type logErr struct {
	Code       int    `json:"code"`
	RemoteAddr string `json:"remoteAddr"`
	ClientID   string `json:"clientID"`
	Error      string `json:"error"`
	Method     string `json:"method"`
	RequestURI string `json:"requestURI"`
}

func (e logErr) writeErrorToLog() {
	log.Printf(
		"%d\t%s\t%s\t%s\t%s\t%s",
		e.Code,
		e.RemoteAddr,
		e.ClientID,
		e.Error,
		e.Method,
		e.RequestURI,
	)
}
