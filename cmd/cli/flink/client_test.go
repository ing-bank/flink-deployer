package flink

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/assert"
)

func TestConstructUrlShouldProperlyFormTheCompleteURL(t *testing.T) {
	api := FlinkRestClient{"http://localhost:80", &retryablehttp.Client{}}

	url := api.constructURL("jobs")

	assert.Equal(t, "http://localhost:80/jobs", url)
}

func createTestServerWithoutBodyCheck(t *testing.T, expectedURL string, status int, body string) *httptest.Server {
	return createTestServer(t, expectedURL, false, "", status, body)
}

func createTestServerWithBodyCheck(t *testing.T, expectedURL string, expectedBody string, status int, body string) *httptest.Server {
	return createTestServer(t, expectedURL, true, expectedBody, status, body)
}

func createTestServer(t *testing.T, expectedURL string, verifyBody bool, expectedBody string, status int, body string) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.URL.String(), expectedURL)

		if verifyBody {
			reqBody, err := ioutil.ReadAll(req.Body)
			assert.Nil(t, err)
			assert.Equal(t, expectedBody, strings.Replace(string(reqBody[:]), "\n", "", -1))
			req.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))
		}

		rw.WriteHeader(status)
		rw.Write([]byte(body))
	}))

	return server
}
