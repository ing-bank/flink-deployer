package operations

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func downloadFile(URL string, apiToken string, targetPath string) (written int64, err error) {
	req, _ := http.NewRequest("GET", URL, nil)
	if len(apiToken) > 0 {
		req.Header.Add("PRIVATE-TOKEN", apiToken)
	}

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return -1, fmt.Errorf("retrieving remote JAR returned unexpected response code: %v", res.StatusCode)
	}

	out, err := os.Create(targetPath)
	if err != nil {
		return
	}
	defer out.Close()

	n, err := io.Copy(out, res.Body)
	if err != nil {
		return n, err
	}

	log.Printf("bytes downloaded: %v", n)
	return n, nil
}
