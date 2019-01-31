package main

import (
	"context"
	"github.com/cfrye2000/productPromisedEventMS/external/github.com/garyburd/redigo/redis"
	"github.com/cfrye2000/productPromisedEventMS/external/github.com/robfig/config"
	"github.com/cfrye2000/productPromisedEventMS/oauth2"
	"log"
	"time"

	"cloud.google.com/go/pubsub"
)

func initialize() {

	//get the config file.
	var e error
	c, e = config.ReadDefault("productPromisedEventMS.cfg")
	if e != nil {
		log.Fatal("config file not found")
	}

	//use port from config file if it exists
	port, _ = c.String("DEFAULT", "port")

	if len(port) == 0 {
		port = "8080"
	}

	//use secure from config file if it exists
	secure, _ = c.Bool("DEFAULT", "ssl")

	if secure {
		port = "443"
	}

	//set up Redis Cache for Oauth2
	redisIP, _ := c.String("oauth2", "redisIPaddress")
	redisPort, _ := c.String("oauth2", "redisPort")

	pool = newPool(redisIP+":"+redisPort, "")

	if err := oauth2.SetRedisTokenCache(pool); err != nil {
		log.Fatal(err)
	}

	//set up for GCS writing
	ctx = context.Background()

	//set up project id
	projectID, _ = c.String("gcp", "projectID")

	// Creates a pubsub client.

	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatal("Failed to create PubSub client: " + err.Error())
	} else {
		psClient = *client
	}

}

func newPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			/*if _, err := c.Do("AUTH", password); err != nil {
			    c.Close()
			    return nil, err
			}*/
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
