package oauth

import (
	"github.com/cfrye2000/productPromisedEventMS/external/github.com/garyburd/redigo/redis"
	"github.com/cfrye2000/productPromisedEventMS/external/github.com/robfig/config"
	"github.com/cfrye2000/productPromisedEventMS/oauth2"
	"testing"
	"time"
)

var testPool *redis.Pool
var poolset bool

func TestBasicAuth(t *testing.T) {
	if _, err := AuthenticateClient("Basic Y2ZyeWU6ZmF0c28="); err != nil {
		t.Errorf("Error authenticating valid client")
	}

	if _, err := AuthenticateClient("Basic Y2ZyeWU6ZG9nc2hpdA=="); err == nil {
		t.Errorf("Error authenticating invalid client")
	}

	if _, err := AuthenticateClient("basic"); err == nil {
		t.Errorf("Error authenticating lowercase missing auth string")
	}
	if _, err := AuthenticateClient("basic Y2ZyeWU6ZmF0c28="); err == nil {
		t.Errorf("Error authenticating lowercase auth string")
	}

	if _, err := AuthenticateClient("Basic x"); err == nil {
		t.Errorf("Error authenticating malformed auth string")
	}

	if _, err := AuthenticateClient("Basic"); err == nil {
		t.Errorf("Error authenticating emtpy auth string")
	}
}

func TestTokenCreation(t *testing.T) {
	//set up
	c, e := config.ReadDefault("../productPromisedEventMS.cfg")
	if e != nil {
		t.Errorf("Error reading config file")
	}

	//set up Redis Cache for Oauth2
	var redisIP, redisPort string
	if s, err := c.String("oauth2", "redisIPaddress"); err != nil {
		t.Errorf(err.Error())
	} else {
		redisIP = s
	}

	if s, err := c.String("oauth2", "redisPort"); err != nil {
		t.Errorf(err.Error())
	} else {
		redisPort = s
	}

	if !poolset {
		testPool = newPool(redisIP+":"+redisPort, "")
		poolset = true
	}

	if err := oauth2.SetRedisTokenCache(testPool); err != nil {
		t.Errorf("Error setting up redis token cache")
	}
	//request a token
	if _, err := oauth2.RequestToken("chris", "client", 1, 1, 1); err != nil {
		t.Errorf("Error requesting oauth2 token: %v", err)
	}
}

func TestTokenRetrieval(t *testing.T) {
	//set up
	c, e := config.ReadDefault("../productPromisedEventMS.cfg")
	if e != nil {
		t.Errorf("Error reading config file")
	}

	//set up Redis Cache for Oauth2
	var redisIP, redisPort string
	if s, err := c.String("oauth2", "redisIPaddress"); err != nil {
		t.Errorf(err.Error())
	} else {
		redisIP = s
	}

	if s, err := c.String("oauth2", "redisPort"); err != nil {
		t.Errorf(err.Error())
	} else {
		redisPort = s
	}

	if !poolset {
		testPool = newPool(redisIP+":"+redisPort, "")
		poolset = true
	}

	if err := oauth2.SetRedisTokenCache(testPool); err != nil {
		t.Errorf("Error setting up redis token cache")
	}

	//request a token
	var token oauth2.TokenResponse
	var accessToken oauth2.AccessToken

	if response, err := oauth2.RequestToken("chris", "client", 1, 1, 5); err != nil {
		t.Errorf("Error getting oauth2 token")
	} else {
		token = response
	}

	//get the token
	if toke, err := oauth2.AuthToken(token.Token); err != nil {
		t.Errorf("Error authenticating oauth2 token: " + err.Error())
	} else {
		accessToken = toke
	}

	if token.Token != accessToken.Token {
		t.Errorf("AccessToken returned doesn't match Token requested")
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
