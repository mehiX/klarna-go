package klarna

import "net/http"

type Option func(c *Client)

func WithDebug(dbg bool) Option {
	return func(c *Client) {
		c.debug = dbg
	}
}

func WithHttpClient(hc http.Client) Option {
	return func(c *Client) {
		c.httpClient = hc
	}
}
