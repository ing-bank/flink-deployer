package flink

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type TerminateJobErrorResponse struct {
	ErrInfo string `json:"error"`
}

// Terminate terminates a running job specified by job ID
func (c FlinkRestClient) Terminate(jobID string, mode string) error {
	var url string
	if len(mode) > 0 {
		url = c.constructURL(fmt.Sprintf("jobs/%v?mode=%v", jobID, mode))
	} else {
		url = c.constructURL(fmt.Sprintf("jobs/%v", jobID))
	}

	req, err := http.NewRequest("PATCH", url, nil)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return err
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	if res.StatusCode != 202 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		return fmt.Errorf("Unexpected response status %v with body %v", res.StatusCode, string(body[:]))
	}

	return nil
}
