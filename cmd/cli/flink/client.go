package flink

import (
	"fmt"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
)

// A FlinkRestClient is a client to interface with
// the Apache Flink REST API
type FlinkRestClient struct {
	BaseURL string
	Client  *retryablehttp.Client
}

func (c FlinkRestClient) constructURL(path string) string {
	return fmt.Sprintf("%v/%v", c.BaseURL, path)
}
