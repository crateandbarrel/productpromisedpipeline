package main

import (
	"context"
	"log"
	"net/http"

	"cloud.google.com/go/pubsub"
	"github.com/cfrye2000/productPromisedEventMS/external/github.com/garyburd/redigo/redis"
	"github.com/cfrye2000/productPromisedEventMS/external/github.com/robfig/config"
)

var c *config.Config
var port string
var secure bool
var pool *redis.Pool
var projectID string
var psClient pubsub.Client
var ctx context.Context

func main() {

	initialize()

	router := NewRouter()
	log.Printf("Listening...%v \n", port)
	if secure {
		key, _ := c.String("ssl", "key")
		cert, _ := c.String("ssl", "cert")
		log.Printf("using...key: %v and cert: %v \n", key, cert)
		log.Fatal(http.ListenAndServeTLS(":"+port, cert, key, router))
	} else {
		log.Fatal(http.ListenAndServe(":"+port, router))
	}
}
