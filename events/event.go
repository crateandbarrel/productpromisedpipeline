package events

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/cfrye2000/productPromisedEventMS/logger"
)

//RootLevel used to for xml
type RootLevel struct {
	XMLName               xml.Name               `xml:"root"`
	ProductPromisedEvents []ProductPromisedEvent `xml:"ProductPromisedEvent"`
}

//ProductPromisedEvent used to model a product event
type ProductPromisedEvent struct {
	EventSource          string    `xml:"EventSource" json:"EventSource"`
	OrderNumber          int       `xml:"OrderNumber" json:"OrderNumber"`
	TransactionNumber    int       `xml:"TransactionNumber" json:"TransactionNumber"`
	CompanyNumber        int       `xml:"CompanyNumber" json:"CompanyNumber"`
	RecipientNumber      int       `xml:"RecipientNumber" json:"RecipientNumber"`
	LineNumber           int       `xml:"LineNumber" json:"LineNumber"`
	SkuNumber            int       `xml:"SkuNumber" json:"SkuNumber"`
	LocationID           int       `xml:"LocationId" json:"LocationId"`
	DmLocationID         int       `xml:"DmLocationId" json:"DmLocationId"`
	SellingSourceCode    int       `xml:"SellingSourceCode" json:"SellingSourceCode"`
	ShippingLocationID   int       `xml:"ShippingLocationId" json:"ShippingLocationId"`
	DemandLocationID     int       `xml:"DemandLocationId" json:"DemandLocationId"`
	PromisedTimestamp    time.Time `xml:"PromisedTimestamp" json:"PromisedTimestamp"`
	OldETATimestamp      time.Time `xml:"OldETATimestamp" json:"OldETATimestamp"`
	NewETATimestamp      time.Time `xml:"NewETATimestamp" json:"NewETATimestamp"`
	TransactionQuantity  int       `xml:"TransactionQuantity" json:"TransactionQuantity"`
	EventCreateTimestamp time.Time `xml:"EventCreateTimestamp" json:"EventCreateTimestamp"`
}

//SubmitEvent  Submit the  Event to the Pubsub service
func SubmitEvent(events []ProductPromisedEvent, psClient pubsub.Client, ctx context.Context) ([]ProductPromisedEvent, error) {

	var returnError error
	var submitEvents []ProductPromisedEvent
	for _, event := range events {
		event.EventCreateTimestamp = time.Now()
		submitEvents = append(submitEvents, event)
	}

	go submitToTopic(submitEvents, psClient, ctx)
	return submitEvents, returnError

}

func submitToTopic(events []ProductPromisedEvent, psClient pubsub.Client, ctx context.Context) {

	for _, event := range events {
		d, err := json.Marshal(event)
		eventString := string(d)
		if err != nil {
			logger.LogEvent("ProductPromisedEvent: "+err.Error(), eventString, "cb-en-us")
		} else {
			topic := psClient.Topic("ProductPromisedEvents")
			_, err = topic.Publish(ctx, &pubsub.Message{Data: d}).Get(ctx)
			if err != nil {
				logger.LogEvent("ProductPromisedEvent: "+err.Error(), eventString, "cb-en-us")
			}
		}
	}

}
