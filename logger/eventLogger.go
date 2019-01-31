package logger

import (
	"log"
	"os"
)

func init() {
	log.SetOutput(os.Stdout)
}

//LogEvent  Used to capture special events for the logs
func LogEvent(eventName string, eventData string, locale string) {

	log.Printf(
		"\t%s\t%s\t%s\t%s",
		"EVT",
		eventName,
		locale,
		eventData,
	)
}
