package flink

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Job struct {
	Id     string `json:"jid"`
	Name   string `json:"name"`
	Status string `json:"state"`
}

type retrieveJobsResponse struct {
	Jobs []Job `json:"jobs"`
}

func (c FlinkRestClient) RetrieveJobs() ([]Job, error) {
	res, err := c.Client.Get(c.constructURL("jobs/overview"))
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
