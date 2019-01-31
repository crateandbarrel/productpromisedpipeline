package oauth2

import (
	"github.com/cfrye2000/productPromisedEventMS/external/github.com/garyburd/redigo/redis"
	"github.com/cfrye2000/productPromisedEventMS/external/github.com/robfig/config"
	"testing"
	"time"
)

var testPool *redis.Pool
var poolset bool

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

	if err := SetRedisTokenCache(testPool); err != nil {
		t.Errorf("Error setting up redis token cache")
	}

	//request a token
	if _, err := RequestToken("chris", "client", 1, 1, 1); err != nil {
		t.Errorf("Error getting oauth2 token")
	}
}

func TestTokenRetrieval(t *testing.T) {
	//set up
	c, e := config.ReadDefault("../productPromisedEventMS.cfg")
	if e != nil {
		t.Errorf("Error reading config file")
	}

	//set up Redis Cache for Oauth2
	redisIP, _ := c.String("oauth2", "redisIPaddress")
	redisPort, _ := c.String("oauth2", "redisPort")

	if !poolset {
		testPool = newPool(redisIP+":"+redisPort, "")
		poolset = true
	}

	if err := SetRedisTokenCache(testPool); err != nil {
		t.Errorf("Error setting up redis token cache")
	}

	//request a token
	var token TokenResponse
	var accessToken AccessToken

	if response, err := RequestToken("chris", "client", 1, 1, 1); err != nil {
		t.Errorf("Error getting oauth2 token")
	} else {
		token = response
	}

	//get the token
	if toke, err := AuthToken(token.Token); err != nil {
		t.Errorf("Error authenticating oauth2 token")
	} else {
		accessToken = toke
	}

	if token.Token != accessToken.Token {
		t.Errorf("AccessToken returned doesn't match Token requested")
	}
}

func TestTokenRefresh(t *testing.T) {
	//set up
	c, e := config.ReadDefault("../productPromisedEventMS.cfg")
	if e != nil {
		t.Errorf("Error reading config file")
	}

	//set up Redis Cache for Oauth2
	redisIP, _ := c.String("oauth2", "redisIPaddress")
	redisPort, _ := c.String("oauth2", "redisPort")

	if !poolset {
		testPool = newPool(redisIP+":"+redisPort, "")
		poolset = true
	}

	if err := SetRedisTokenCache(testPool); err != nil {
		t.Errorf("Error setting up redis token cache")
	}

	//request a token
	var token TokenResponse
	var accessToken AccessToken

	if response, err := RequestToken("chris", "client", 1, 1, 1); err != nil {
		t.Errorf("Error getting oauth2 token")
	} else {
		token = response
	}

	//get the token
	if toke, err := AuthToken(token.Token); err != nil {
		t.Errorf("Error authenticating oauth2 token")
	} else {
		accessToken = toke
	}

	if token.Token != accessToken.Token {
		t.Errorf("AccessToken refreshed doesn't match Token requested")
	}

	//refresh the token
	var refreshedToken AccessToken
	if _, err := RefreshToken("chris", accessToken.Token, "refresh_token", 1, 1, 2); err != nil {
		t.Errorf("Error refreshing oauth2 token")
	}

	//get the token
	if toke, err := AuthToken(token.Token); err != nil {
		t.Errorf("Error authenticating oauth2 token: %v", err.Error())
	} else {
		refreshedToken = toke
	}

	if refreshedToken.Expiration-accessToken.Expiration != 1 {
		t.Errorf("AccessToken refreshed didn't work")
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
