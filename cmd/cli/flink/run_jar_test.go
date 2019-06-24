package flink

import (
	"net/http"
	"testing"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/assert"
)

func TestRunJarReturnsAnErrorWhenTheStatusIsNot200(t *testing.T) {
	server := createTestServerWithBodyCheck(t, "/jars/id/run", `{"entryClass":"MainClass","programArgs":"","parallelism":1,"allowNonRestoredState":false,"savepointPath":"/data/flink"}`, http.StatusAccepted, "{}")
	defer server.Close()

	api := FlinkRestClient{
		BaseURL: server.URL,
		Client:  retryablehttp.NewClient(),
	}
	err := api.RunJar("id", "MainClass", []string{}, 1, "/data/flink", false)

	assert.EqualError(t, err, "Unexpected response status 202 with body {}")
}

func TestRunJarCorrectlyReturnsNilWhenTheCallSucceeds(t *testing.T) {
	server := createTestServerWithBodyCheck(t, "/jars/id/run", `{"entryClass":"MainClass","programArgs":"","parallelism":1,"allowNonRestoredState":false,"savepointPath":"/data/flink"}`, http.StatusOK, "")
	defer server.Close()

	api := FlinkRestClient{
		BaseURL: server.URL,
		Client:  retryablehttp.NewClient(),
	}
	err := api.RunJar("id", "MainClass", []string{}, 1, "/data/flink", false)

	assert.Nil(t, err)
}
