package flink

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/assert"
)

func TestConstructUrlShouldProperlyFormTheCompleteURL(t *testing.T) {
	api := FlinkRestClient{
		BaseURL: "http://localhost:80",
		Client:  &retryablehttp.Client{},
	}

	url := api.constructURL("jobs")

	assert.Equal(t, "http://localhost:80/jobs", url)
}

func TestNewRequestShouldAddTheBasicAuthenticationHeadersWhenTheCredentialsAreSet(t *testing.T) {
	api := FlinkRestClient{"http://localhost:80", "username", "password", &retryablehttp.Client{}}

	req, err := api.newRequest("GET", "jobs", nil)

	assert.Nil(t, err)
	assert.Equal(t, "Basic dXNlcm5hbWU6cGFzc3dvcmQ=", req.Header.Get("Authorization"))
}

func TestNewRequestShouldNotAddTheBasicAuthenticationHeadersWhenTheUsernameIsUnset(t *testing.T) {
	api := FlinkRestClient{
		BaseURL:           "http://localhost:80",
		BasicAuthPassword: "password",
		Client:            &retryablehttp.Client{},
	}

	req, err := api.newRequest("GET", "jobs", nil)

	assert.Nil(t, err)
	assert.Equal(t, "", req.Header.Get("Authorization"))
}

func TestNewRequestShouldNotAddTheBasicAuthenticationHeadersWhenThePassworsdIsUnset(t *testing.T) {
	api := FlinkRestClient{
		BaseURL:           "http://localhost:80",
		BasicAuthUsername: "username",
		Client:            &retryablehttp.Client{},
	}

	req, err := api.newRequest("GET", "jobs", nil)

	assert.Nil(t, err)
	assert.Equal(t, "", req.Header.Get("Authorization"))
}

func TestNewRequestShouldNotAddTheBasicAuthenticationHeadersWhenBothCredentialsAreUnset(t *testing.T) {
	api := FlinkRestClient{
		BaseURL: "http://localhost:80",
		Client:  &retryablehttp.Client{},
	}

	req, err := api.newRequest("GET", "jobs", nil)

	assert.Nil(t, err)
	assert.Equal(t, "", req.Header.Get("Authorization"))
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
