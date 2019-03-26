package flink

import (
	"net/http"
	"testing"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/assert"
)

func TestRetrieveJobsReturnsAnErrorWhenTheStatusIsNot200(t *testing.T) {
	server := createTestServerWithBodyCheck(t, "/jobs/overview", "", http.StatusAccepted, "{}")
	defer server.Close()

	api := FlinkRestClient{
		BaseURL: server.URL,
		Client:  retryablehttp.NewClient(),
	}
	_, err := api.RetrieveJobs()

	assert.EqualError(t, err, "Unexpected response status 202 with body {}")
}

func TestRetrieveJobsReturnsAnErrorWhenItCannotDeserializeTheResponseAsJSON(t *testing.T) {
	server := createTestServerWithBodyCheck(t, "/jobs/overview", "", http.StatusOK, `{"jobs: []}`)
	defer server.Close()

	api := FlinkRestClient{
		BaseURL: server.URL,
		Client:  retryablehttp.NewClient(),
	}
	_, err := api.RetrieveJobs()

	assert.EqualError(t, err, "Unable to parse API response as valid JSON: {\"jobs: []}")
}

func TestRetrieveJobsCorrectlyReturnsAnArrayOfJobs(t *testing.T) {
	server := createTestServerWithBodyCheck(t, "/jobs/overview", "", http.StatusOK, `{"jobs":[{"jid": "1", "name": "Job A", "state": "RUNNING"}]}`)
	defer server.Close()

	api := FlinkRestClient{
		BaseURL: server.URL,
		Client:  retryablehttp.NewClient(),
	}
	jobs, err := api.RetrieveJobs()

	assert.Len(t, jobs, 1)
	assert.Nil(t, err)
}
