package flink

import (
	"net/http"
	"testing"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/assert"
)

/*
 * Create Savepoint
 */
func TestCreateSavepointReturnsAnErrorWhenTheStatusIsNot202(t *testing.T) {
	server := createTestServerWithBodyCheck(t, "/jobs/1/savepoints", `{"target-directory":"/data/flink","cancel-job":false}`, http.StatusOK, "{}")
	defer server.Close()

	api := FlinkRestClient{
		BaseURL: server.URL,
		Client:  retryablehttp.NewClient(),
	}
	_, err := api.CreateSavepoint("1", "/data/flink")

	assert.EqualError(t, err, "Unexpected response status 200 with body {}")
}

func TestCreateSavepointReturnsAnErrorWhenItCannotDeserializeTheResponseAsJSON(t *testing.T) {
	server := createTestServerWithBodyCheck(t, "/jobs/1/savepoints", `{"target-directory":"/data/flink","cancel-job":false}`, http.StatusAccepted, `{"jobs: []}`)
	defer server.Close()

	api := FlinkRestClient{
		BaseURL: server.URL,
		Client:  retryablehttp.NewClient(),
	}
	_, err := api.CreateSavepoint("1", "/data/flink")

	assert.EqualError(t, err, "Unable to parse API response as valid JSON: {\"jobs: []}")
}

func TestCreateSavepointCorrectlyReturnsARequestID(t *testing.T) {
	server := createTestServerWithBodyCheck(t, "/jobs/1/savepoints", `{"target-directory":"/data/flink","cancel-job":false}`, http.StatusAccepted, `{"request-id": "1"}`)
	defer server.Close()

	api := FlinkRestClient{
		BaseURL: server.URL,
		Client:  retryablehttp.NewClient(),
	}
	res, err := api.CreateSavepoint("1", "/data/flink")

	assert.Equal(t, res.RequestID, "1")
	assert.Nil(t, err)
}

/*
 * Monitor Savepoint Creation
 */
func TestMonitorSavepointCreationReturnsAnErrorWhenTheStatusIsNot200(t *testing.T) {
	server := createTestServerWithBodyCheck(t, "/jobs/id-1/savepoints/request-id-1", "", http.StatusAccepted, "{}")
	defer server.Close()

	api := FlinkRestClient{
		BaseURL: server.URL,
		Client:  retryablehttp.NewClient(),
	}
	_, err := api.MonitorSavepointCreation("id-1", "request-id-1")

	assert.EqualError(t, err, "Unexpected response status 202 with body {}")
}

func TestMonitorSavepointCreationReturnsAnErrorWhenItCannotDeserializeTheResponseAsJSON(t *testing.T) {
	server := createTestServerWithBodyCheck(t, "/jobs/id-1/savepoints/request-id-1", "", http.StatusOK, `{"jobs: []}`)
	defer server.Close()

	api := FlinkRestClient{
		BaseURL: server.URL,
		Client:  retryablehttp.NewClient(),
	}
	_, err := api.MonitorSavepointCreation("id-1", "request-id-1")

	assert.EqualError(t, err, "Unable to parse API response as valid JSON: {\"jobs: []}")
}

func TestMonitorSavepointCreationCorrectlyReturnsARequestID(t *testing.T) {
	server := createTestServerWithBodyCheck(t, "/jobs/id-1/savepoints/request-id-1", "", http.StatusOK, `{"status":{"id":"PENDING"}}`)
	defer server.Close()

	api := FlinkRestClient{
		BaseURL: server.URL,
		Client:  retryablehttp.NewClient(),
	}
	res, err := api.MonitorSavepointCreation("id-1", "request-id-1")

	assert.Equal(t, res.Status.Id, "PENDING")
	assert.Nil(t, err)
}
