package httpclient

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type HttpClient struct {
	client *resty.Client
}

func (c *HttpClient) New() {
	c.client = resty.New()
}

func (c *HttpClient) Get(url string) (string, error) {
	response, err := c.client.R().Get(url)
	if err != nil {
		return "", fmt.Errorf(`failed to get response, error : %v`, err)
	}
	return string(response.Body()), nil
}

