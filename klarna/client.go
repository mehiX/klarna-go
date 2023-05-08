package klarna

import "net/http"

type Client struct {
	token      string
	targetURL  string
	httpClient http.Client
	debug      bool
}

func New(url string, tkn string, opts ...Option) *Client {
	c := &Client{targetURL: url, token: tkn, debug: false}

	for _, o := range opts {
		o(c)
	}

	return c
}

func (c *Client) Apply(o Option) {
	o(c)
}
