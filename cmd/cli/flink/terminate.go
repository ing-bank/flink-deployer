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
	var path string
	if len(mode) > 0 {
		path = fmt.Sprintf("jobs/%v?mode=%v", jobID, mode)
	} else {
		path = fmt.Sprintf("jobs/%v", jobID)
	}

	c.Client.CheckRetry = RetryPolicy
	req, err := c.newRequest("PATCH", c.constructURL(path), nil)
	res, err := c.Client.Do(req)

	defer res.Body.Close()

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

// Do not retry when status code is 500. (indicating the job is not stoppable)
func RetryPolicy(ctx context.Context, resp *http.Response, err error) (bool, error) {
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
