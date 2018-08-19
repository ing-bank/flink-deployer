package flink

import (
	"fmt"
)

func (c FlinkRestClient) Cancel(jobID string) error {
	res, err := c.Client.Get(c.constructURL(fmt.Sprintf("jobs/%v/yarn-cancel", jobID)))
	if err != nil {
		return err
	}

	if res.StatusCode != 202 {
		return fmt.Errorf("Unexpected response status %v", res.StatusCode)
	}

	return nil
}
