package client

import "net/url"

// base struct for all web clients
type webClient struct {
	url *url.URL
}
