package oauth2

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cfrye2000/productPromisedEventMS/external/code.google.com/p/go-uuid/uuid"
	"github.com/cfrye2000/productPromisedEventMS/external/github.com/garyburd/redigo/redis"
)

var conn redis.Conn
var connectionSet = false

//5 day buffer = 5 days * 24 hours/day * 60 minutes/hour * 60 seconds/minute = 432000 seconds
const redisExpiryBuffer = 432000

//TokenResponse Used to model a newly created token
type TokenResponse struct {
	Token     string `json:"access_token"`
	TokenType string `json:"token_type"`
	ExpiresIn string `json:"expires_in"`
}

//AccessToken Used to model a token from the data store
type AccessToken struct {
	Token         string `json:"access_token"`
	ClientID      string `json:"client_id"`
	GrantType     string `json:"grant_type"`
	Expiration    int64  `json:"expiration"`
	SecurityLevel int    `json:"secruity_level"`
	HitsPerMinute int64  `json:"hits_per_minute"`
}

//AuthResponse Used to model authorization request for a token
type AuthResponse struct {
	Vali2d    bool   `json:"is_valid"`
	Token     string `json:"access_token"`
	TokenType string `json:"token_type"`
	ExpiresIn string `json:"expires_in"`
}

var (
	pool *redis.Pool
)

func init() {

}

//SetRedisTokenCache Used to provide access to a particular Redis Store for tokens
func SetRedisTokenCache(p *redis.Pool) error {
	if !connectionSet {
		pool = p
		connectionSet = true
	}

	return nil
}

//AuthToken Used to retreive a token from the data store
func AuthToken(tokenString string) (AccessToken, error) {
	var token AccessToken
	var returnError error

	if !connectionSet {
		returnError = errors.New("No redis token cache established")
	} else {
		conn := pool.Get()
		defer conn.Close()
		tokenItem, err := redis.Bytes(conn.Do("GET", getAccessTokenKey(tokenString)))
		if err != nil {
			returnError = errors.New("token not found or expired")
		} else {
			var t AccessToken
			err := json.Unmarshal(tokenItem, &t)
			if err != nil {
				returnError = errors.New("token not found or expired")
			} else {
				token = t
				//check time to live
				if ttl, err := redis.Int(conn.Do("TTL", getAccessTokenKey(tokenString))); err != nil || (ttl > 0 && ttl < redisExpiryBuffer) {
					returnError = errors.New("token is expired, please refresh")
				}
			}
		}
	}

	return token, returnError
}

//RequestToken Used to create a new Token
func RequestToken(clientID string, grantType string, securityLevel int, hitsPerMinute int64, expiry int64) (TokenResponse, error) {
	var tokenResponse TokenResponse
	var returnError error

	if !connectionSet {
		returnError = errors.New("No redis token cache established")
	} else {
		conn := pool.Get()
		defer conn.Close()
		uuid := uuid.NewUUID()
		//add a buffer to leave it in Redis after it has expired so they have time to refresh
		expirationTime := expiry + redisExpiryBuffer
		accessToken := AccessToken{uuid.String(), clientID, grantType, expirationTime, securityLevel, hitsPerMinute}
		if bytes, err := json.Marshal(accessToken); err != nil {
			returnError = err
		} else {
			// Add the item to redis, if the key does not already exist
			_, err := redis.Bytes(conn.Do("SET", getAccessTokenKey(uuid.String()), bytes))
			if err != nil {
				returnError = err
			} else {
				if _, err := redis.Int(conn.Do("EXPIRE", getAccessTokenKey(uuid.String()), expirationTime)); err != nil {
					returnError = err
				}
			}
		}

		if returnError == nil {
			tokenResponse = TokenResponse{uuid.String(), "Bearer", fmt.Sprintf("%v", expiry)}
		}
	}

	return tokenResponse, returnError
}

//RefreshToken Used to create a new Token
func RefreshToken(clientID string, tokenString string, grantType string, securityLevel int, hitsPerMinute int64, expiry int64) (TokenResponse, error) {
	var tokenResponse TokenResponse
	var returnError error

	if !connectionSet {
		returnError = errors.New("No redis token cache established")
	} else {
		var token AccessToken
		conn := pool.Get()
		defer conn.Close()
		tokenItem, err := redis.Bytes(conn.Do("GET", getAccessTokenKey(tokenString)))
		if err != nil {
			returnError = errors.New("token not found or long expired")
		} else {
			var t AccessToken
			err := json.Unmarshal(tokenItem, &t)
			if err != nil {
				returnError = errors.New("token not found or long expired")
			} else {
				token = t
				//add a buffer to leave it in Redis after it has expired so they have time to refresh
				expirationTime := expiry + redisExpiryBuffer
				accessToken := AccessToken{token.Token, clientID, grantType, expirationTime, securityLevel, hitsPerMinute}
				if bytes, err := json.Marshal(accessToken); err != nil {
					returnError = err
				} else {
					// Add the item to redis, if the key does not already exist
					_, err := redis.Bytes(conn.Do("SET", getAccessTokenKey(token.Token), bytes))
					if err != nil {
						returnError = err
					} else {
						if _, err := redis.Int(conn.Do("EXPIRE", getAccessTokenKey(token.Token), expirationTime)); err != nil {
							returnError = err
						}
					}
				}
			}

			if returnError == nil {
				tokenResponse = TokenResponse{token.Token, "Bearer", fmt.Sprintf("%v", expiry)}
			}
		}
	}

	return tokenResponse, returnError
}

func getAccessTokenKey(token string) string {
	return token + ":OAuthAccessToken"
}
