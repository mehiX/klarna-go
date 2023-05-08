package klarna

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
)

func (c *Client) Post(ctx context.Context, endpoint string, body []byte) ([]byte, error) {

	url := c.targetURL + endpoint

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("[Post] creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Token "+c.token)

	if c.debug {
		b, _ := httputil.DumpRequestOut(req, true)
		fmt.Println(string(b))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if c.debug {
		b, _ := httputil.DumpResponse(resp, true)
		fmt.Printf("\n%v\n\n", string(b))
	}

	return io.ReadAll(resp.Body)
}

func (c *Client) Get(ctx context.Context, endpoint string) ([]byte, error) {

	url := c.targetURL + endpoint

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Token "+c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
