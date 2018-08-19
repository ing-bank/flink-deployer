package flink

import (
	"fmt"
	"net/http"
)

type FlinkRestClient struct {
	BaseURL string
	Client  *http.Client
}

func (c FlinkRestClient) constructURL(path string) string {
	return fmt.Sprintf("%v/%v", c.BaseURL, path)
}
