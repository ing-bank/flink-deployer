package flink

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// A Job is a representation for a Flink Job
type Job struct {
	ID     string `json:"jid"`
	Name   string `json:"name"`
	Status string `json:"state"`
}

type retrieveJobsResponse struct {
	Jobs []Job `json:"jobs"`
}

// RetrieveJobs returns all the jobs on the Flink cluster
func (c FlinkRestClient) RetrieveJobs() ([]Job, error) {
	req, err := c.newRequest("GET", c.constructURL("jobs/overview"), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []Job{}, err
	}

	if res.StatusCode != 200 {
		return []Job{}, fmt.Errorf("Unexpected response status %v with body %v", res.StatusCode, string(body[:]))
	}

	retrieveJobsResponse := retrieveJobsResponse{}
	err = json.Unmarshal(body, &retrieveJobsResponse)
	if err != nil {
		return []Job{}, fmt.Errorf("Unable to parse API response as valid JSON: %v", string(body[:]))
	}

	return retrieveJobsResponse.Jobs, nil
}
