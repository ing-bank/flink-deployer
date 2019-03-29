package flink

import (
	"context"
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

	req, err := c.newRequest("PATCH", url, nil)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return err
	}

	// set retry policy: do not retry when status code is 500. (indicating the job is not stoppable)
	c.Client.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (b bool, e error) {
		if ctx.Err() != nil {
			return false, ctx.Err()
		}

		if err != nil {
			return true, err
		}

		if resp.StatusCode == 0 || resp.StatusCode > 500 {
			return true, nil
		}

		return false, nil
	}

	res, err := c.Client.Do(req)

	if err != nil {
		return err
	}

	if res.StatusCode != 202 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		return fmt.Errorf("Unexpected response status %v with body %v", res.StatusCode, string(body[:]))
	}

	return nil
}
