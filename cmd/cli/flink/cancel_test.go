package flink

import (
	"net/http"
	"testing"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/assert"
)

func TestCancelReturnsAnErrorWhenTheResponseStatusIsNot202(t *testing.T) {
	server := createTestServerWithBodyCheck(t, "/jobs/id", "", http.StatusOK, "OK")
	defer server.Close()

	api := FlinkRestClient{
		BaseURL: server.URL,
		Client:  retryablehttp.NewClient(),
	}
	err := api.Cancel("id")

	assert.EqualError(t, err, "Unexpected response status 200")
}

func TestCancelShouldNotReturnAnErrorWhenTheResponseStatusIs202(t *testing.T) {
	server := createTestServerWithBodyCheck(t, "/jobs/id", "", http.StatusAccepted, "")
	defer server.Close()

	api := FlinkRestClient{
		BaseURL: server.URL,
		Client:  retryablehttp.NewClient(),
	}
	err := api.Cancel("id")

	assert.Nil(t, err)
}
