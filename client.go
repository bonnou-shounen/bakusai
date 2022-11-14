package bakusai

import (
	"context"
	"fmt"
	"net/http"
)

type Client struct {
	HTTPClient *http.Client
}

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Client{
		HTTPClient: httpClient,
	}
}

func (c *Client) Get(ctx context.Context, uri string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf(`on http.NewRequestWithContext(.."%s"..): %w`, uri, err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf(`on http.Client.Do(): %w`, err)
	}

	return resp, nil
}
