package operations

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func callDownloadFile(t *testing.T, apiToken string, targetPath string) (written int64, err error) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("TESTTHIS"))
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got ‘%s’", r.Method)
		}
		if r.Header.Get("PRIVATE-TOKEN") != "header" {
			t.Errorf("Expected certain header, got ‘%s’", r.Header)
		}
	}))
	defer ts.Close()

	return downloadFile(ts.URL, apiToken, targetPath)
}

func TestDownloadFile(t *testing.T) {
	targetPath := "./test.file"
	apiToken := "header"

	callDownloadFile(t, apiToken, targetPath)

	f, _ := ioutil.ReadFile(targetPath)
	assert.Equal(t, "TESTTHIS", string(f))
	os.Remove(targetPath)
}
