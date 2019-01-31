package events

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/cfrye2000/productPromisedEventMS/external/github.com/robfig/config"
	"testing"
)

func TestSubmitEvent(t *testing.T) {
	//set up
	c, e := config.ReadDefault("../productPromisedEventMS.cfg")
	if e != nil {
		t.Errorf("Error reading config file")
	}

	//set up for GCS writing
	ctx := context.Background()

	// Sets your Google Cloud Platform project ID.
	projectID, _ := c.String("gcp", "projectID")

	// Creates a client.
	psClient, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		t.Errorf("Failed to create pubsub client: " + err.Error())
	}

	var events []ProductPromisedEvent
	var event ProductPromisedEvent
	event.OrderNumber = 1
	events = append(events, event)

	if _, err := SubmitEvent(events, *psClient, ctx); err != nil {
		t.Errorf("Error submitting event " + err.Error())
	}
}
