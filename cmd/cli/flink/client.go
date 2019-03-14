package flink

import (
	"fmt"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
)

// A FlinkRestClient is a client to interface with
// the Apache Flink REST API
type FlinkRestClient struct {
	BaseURL           string
	BasicAuthUsername string
	BasicAuthPassword string
	Client            *retryablehttp.Client
}

func (c FlinkRestClient) constructURL(path string) string {
	return fmt.Sprintf("%v/%v", c.BaseURL, path)
}

func basicAuthenticationCredentialsDefined(username, password string) bool {
	return len(username) != 0 && len(password) != 0
}

func (c FlinkRestClient) newRequest(method, url string, rawBody interface{}) (*retryablehttp.Request, error) {
	req, err := retryablehttp.NewRequest(method, url, rawBody)
	if err != nil {
		return nil, err
	}

	if basicAuthenticationCredentialsDefined(c.BasicAuthUsername, c.BasicAuthPassword) == true {
		req.SetBasicAuth(c.BasicAuthUsername, c.BasicAuthPassword)
	}

	return req, err
}
