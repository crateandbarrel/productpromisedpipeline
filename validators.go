package main

import (
	"github.com/cfrye2000/productPromisedEventMS/events"
)

type validationFunc func(event events.ProductPromisedEvent) (bool, string)

//Validator structure of a validator
type Validator struct {
	EventName    string
	ValidateFunc validationFunc
}

//Validators list of all the validators
type Validators []Validator

var validators Validators
