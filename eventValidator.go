package main

import (
	"time"

	"github.com/cfrye2000/productPromisedEventMS/events"
)

func init() {

	validator := Validator{
		"ProductPromisedEvent",
		ValidateEvent,
	}

	validators = append(validators, validator)

}

//ValidateEvent  Validate the  Event
func ValidateEvent(event events.ProductPromisedEvent) (bool, string) {

	if len(event.EventSource) == 0 {
		return false, "Specify Event Source"
	}
	if event.OrderNumber == 0 {
		return false, "Specify Order Number"
	}
	if event.TransactionNumber == 0 {
		return false, "Specify Transaction Number"
	}
	if event.CompanyNumber == 0 {
		return false, "Specify Company Number"
	}
	if event.RecipientNumber == 0 {
		return false, "Specify Recipient Number"
	}
	if event.LineNumber == 0 {
		return false, "Specify Line Number"
	}
	if event.SkuNumber == 0 {
		return false, "Specify Sku Number"
	}
	if event.DmLocationID == 0 {
		return false, "Specify Dm Location Id"
	}
	if event.LocationID == 0 {
		return false, "Specify Location ID"
	}
	if event.TransactionQuantity == 0 {
		return false, "Specify Transaction Quantity"
	}

	notime, _ := time.Parse(time.RFC3339, "0001-01-01T00:00:00Z")

	if event.PromisedTimestamp == notime {
		return false, "Specify Promised Timestamp"
	}

	if event.OldETATimestamp == notime {
		return false, "Specify Old ETA Timestamp"
	}

	if event.NewETATimestamp == notime {
		return false, "Specify New ETA Timestamp"
	}

	return true, ""

}
