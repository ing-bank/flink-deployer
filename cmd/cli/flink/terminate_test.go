package flink

import (
	"net/http"
	"testing"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/assert"
)

func TestTerminateWithModeCancelAndStatusSuccess(t *testing.T) {
	server := createTestServerWithBodyCheck(t, "/jobs/id?mode=cancel", "", http.StatusAccepted, "")
	defer server.Close()

	api := FlinkRestClient{server.URL, retryablehttp.NewClient()}
	err := api.Terminate("id", "cancel")

	assert.Nil(t, err)
}

func TestTerminateWithModeCancelAndStatus404(t *testing.T) {
	server := createTestServerWithBodyCheck(t, "/jobs/id?mode=cancel", "", http.StatusNotFound, "not found")
	defer server.Close()

	api := FlinkRestClient{server.URL, retryablehttp.NewClient()}
	err := api.Terminate("id", "cancel")

	assert.EqualError(t, err, "Unexpected response status 404 with body not found")
}

func TestTerminateWithModeStopAndStatusSuccess(t *testing.T) {
	server := createTestServerWithBodyCheck(t, "/jobs/id?mode=stop", "", http.StatusAccepted, "OK")
	defer server.Close()

	api := FlinkRestClient{server.URL, retryablehttp.NewClient()}
	err := api.Terminate("id", "stop")

	assert.Nil(t, err)
}

func TestTerminateWithModeStopAndStatusFailure(t *testing.T) {
	server := createTestServerWithBodyCheck(t, "/jobs/id?mode=stop", "", http.StatusInternalServerError, "error")
	defer server.Close()

	api := FlinkRestClient{server.URL, retryablehttp.NewClient()}
	err := api.Terminate("id", "stop")

	assert.EqualError(t, err, "Unexpected response status 500 with body error")
}
