package flink

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCancelReturnsAnErrorWhenTheResponseStatusIsNot202(t *testing.T) {
	server := createTestServerWithBodyCheck(t, "/jobs/id", "", http.StatusOK, "OK")
	defer server.Close()

	api := FlinkRestClient{server.URL, server.Client()}
	err := api.Cancel("id")

	assert.EqualError(t, err, "Unexpected response status 200")
}

func TestCancelShouldNotReturnAnErrorWhenTheResponseStatusIs202(t *testing.T) {
	server := createTestServerWithBodyCheck(t, "/jobs/id", "", http.StatusAccepted, "")
	defer server.Close()

	api := FlinkRestClient{server.URL, server.Client()}
	err := api.Cancel("id")

	assert.Nil(t, err)
}
