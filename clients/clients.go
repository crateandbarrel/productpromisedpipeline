package clients

//Client Is a struct that models authorized users of the API
type Client struct {
	ClientID      string `json:"client_id"`
	Secret        string `json:"secret"`
	SecurityLevel int    `json:"security_level"`
	HitsPerMinute int64  `json:"hits_per_minute"`
	Expiry        int64  `json:"expires_in"`
	Description   string `json:"description"`
	ContactInfo   string `json:"contact_info"`
	AddedBy       string `json:"added_by"`
	Active        bool   `json:"active"`
}

//expires in 29 days = 29 days * 24 hours/day * 60 minutes/hour * 60 seconds/minute = 2592000 seconds
//anything greater than 29 days is taken as a unix timestamp... switch to that after 30 days

//HitsPerMinute of -1 means no limit
//expiry of -1 means never expire

//Clients A collection of authorized users of the API
var Clients = map[string]Client{
	"cfrye":            Client{"cfrye", "fatso", 9, 10, 1296000, "test client", "cfrye2000@gmail.com", "cfrye2000@gmail.com", true},
	"chrislong":        Client{"chrislong", "doggie", 9, -1, 2592000, "test client", "cfrye2000@gmail.com", "cfrye2000@gmail.com", true},
	"chrisimmediate":   Client{"chrisimmediate", "immediate", 9, 10, 30, "immediate expiry", "cfrye2000@gmail.com", "cfrye2000@gmail.com", true},
	"chrisinactive":    Client{"chrisinactive", "inactive", 9, 10, 1296000, "immediate expiry", "cfrye2000@gmail.com", "cfrye2000@gmail.com", false},
	"cfrye2000":        Client{"cfrye2000", "bZ4YK6cFzm2xjXxN", 9, -1, 3155692600, "cfrye2000 client that lasts a long time", "cfrye2000@gmail.com", "cfrye2000@gmail.com", true},
	"microsoftpoc":     Client{"microsoftpoc", "015fa4f21046490aa7ecd4360904e5e0", 9, 100, 1296000, "client used for microsoftPOC", "nstroop@uxceclipse.com", "cfrye2000@gmail.com", true},
	"webserviceorders": Client{"webserviceorders", "173c915cfde5472193abfac700d92c5d", 9, -1, 86400, "client used for order webservice", "bkarger@crateandbarrel.com", "cfrye@crateandbarrel.com", true},
}
