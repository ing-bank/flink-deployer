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
func TestUploadJarReturnsAnErrorWhenTheStatusIsNot200(t *testing.T) {
	server := createTestServerWithoutBodyCheck(t, "/jars/upload", http.StatusAccepted, "{}")
	defer server.Close()

	api := FlinkRestClient{
		BaseURL: server.URL,
		Client:  retryablehttp.NewClient(),
	}
	_, err := api.UploadJar("../testdata/sample.jar")

	assert.EqualError(t, err, "Unexpected response status 202 with body {}")
}

func TestUploadJarReturnsAnErrorWhenItCannotDeserializeTheResponseAsJSON(t *testing.T) {
	server := createTestServerWithoutBodyCheck(t, "/jars/upload", http.StatusOK, `{"jobs: []}`)
	defer server.Close()

	api := FlinkRestClient{
		BaseURL: server.URL,
		Client:  retryablehttp.NewClient(),
	}
	_, err := api.UploadJar("../testdata/sample.jar")

	assert.EqualError(t, err, "Unable to parse API response as valid JSON: {\"jobs: []}")
}

func TestUploadJarCorrectlyReturnsARequestID(t *testing.T) {
	server := createTestServerWithoutBodyCheck(t, "/jars/upload", http.StatusOK, `{"filename": "/flink/jars/sample.jar", "status": "success"}`)
	defer server.Close()

	api := FlinkRestClient{
		BaseURL: server.URL,
		Client:  retryablehttp.NewClient(),
	}
	res, err := api.UploadJar("../testdata/sample.jar")

	assert.Equal(t, res.Filename, "/flink/jars/sample.jar")
	assert.Equal(t, res.Status, "success")
	assert.Nil(t, err)
}
