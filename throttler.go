package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/cfrye2000/productPromisedEventMS/external/github.com/garyburd/redigo/redis"
	"github.com/cfrye2000/productPromisedEventMS/external/github.com/gorilla/context"
	"github.com/cfrye2000/productPromisedEventMS/oauth2"
)

const statusTooManyRequests int = 429

//Throttle struct for persisted throttle data
type Throttle struct {
	ClientID      string `json:"client_id"`
	Expiration    int64  `json:"expiration"`
	HitsRemaining int64  `json:"hits_remaining"`
}

//Throttler  Router that wraps other routers for throttling access by ClientID
func Throttler(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var throttleID string
		var token oauth2.AccessToken

		if val, ok := context.GetOk(r, AccessToken); ok {
			token = val.(oauth2.AccessToken)
			throttleID = token.ClientID
		} else {
			//no token use IP address for throttleID
			token = oauth2.AccessToken{}
			throttleID = r.RemoteAddr
		}

		if token.HitsPerMinute < 0 {
			//this client has been set up to not be throttled
			inner.ServeHTTP(w, r)
		} else {
			var throttle = Throttle{}

			conn := pool.Get()
			defer conn.Close()

			throttleString, err := redis.Bytes(conn.Do("GET", getThrottleKey(throttleID)))
			if err == nil {
				err := json.Unmarshal(throttleString, &throttle)
				if err != nil {
					throttle = Throttle{}
				}
			}

			if throttle.Expiration == 0 {
				//no current throttle stored
				newThrottle(throttleID, token, &throttle)
			}

			expiry := throttle.Expiration
			now := time.Now().Unix()
			timeRemaining := expiry - now
			if timeRemaining <= 0 {
				newThrottle(throttleID, token, &throttle)
				timeRemaining = 60
			}

			throttle.HitsRemaining = throttle.HitsRemaining - 1

			w.Header().Set("X-Rate-Limit-Limit", strconv.FormatInt(token.HitsPerMinute, 10))
			w.Header().Set("X-Rate-Limit-Remaining", strconv.FormatInt(throttle.HitsRemaining, 10))
			w.Header().Set("X-Rate-Limit-Reset", strconv.FormatInt(timeRemaining, 10))

			if bytes, err := json.Marshal(throttle); err != nil {
				logErr{Code: http.StatusInternalServerError, RemoteAddr: r.RemoteAddr, ClientID: throttleID, Error: throttleID + " marshalling error", Method: r.Method, RequestURI: r.RequestURI}.writeErrorToLog()
			} else {
				_, err := redis.Bytes(conn.Do("SET", getThrottleKey(throttleID), bytes))
				if err != nil {
					logErr{Code: http.StatusInternalServerError, RemoteAddr: r.RemoteAddr, ClientID: throttleID, Error: "throttler for " + throttleID + " : " + err.Error() + " : " + strconv.FormatInt(timeRemaining, 10), Method: r.Method, RequestURI: r.RequestURI}.writeErrorToLog()

				} else {
					if _, err := redis.Int(conn.Do("EXPIRE", getThrottleKey(throttleID), timeRemaining)); err != nil {
						logErr{Code: http.StatusInternalServerError, RemoteAddr: r.RemoteAddr, ClientID: throttleID, Error: "throttler for " + throttleID + " : " + err.Error() + " : " + strconv.FormatInt(timeRemaining, 10), Method: r.Method, RequestURI: r.RequestURI}.writeErrorToLog()
					}
				}
			}

			if throttle.HitsRemaining < 0 {
				logErr{Code: statusTooManyRequests, RemoteAddr: r.RemoteAddr, ClientID: throttleID, Error: throttleID + " Exceeded Maximum Hits per Minute", Method: r.Method, RequestURI: r.RequestURI}.writeErrorToLog()
				WriteResponse(w, r, statusTooManyRequests, errorResponse{Code: statusTooManyRequests, Text: throttleID + " Exceeded Maximum Hits per Minute"})
			} else {
				inner.ServeHTTP(w, r)
			}
		}
	})
}

func getThrottleKey(throttleID string) string {
	return throttleID + ":ThrottleData"
}

func newThrottle(throttleID string, token oauth2.AccessToken, throttle *Throttle) {
	throttle.ClientID = throttleID
	throttle.Expiration = time.Now().Unix() + 60
	throttle.HitsRemaining = token.HitsPerMinute
}
