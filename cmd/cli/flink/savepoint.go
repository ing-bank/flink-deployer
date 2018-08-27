package flink

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type createSavepointRequest struct {
	TargetDirectory string `json:"target-directory"`
	CancelJob       bool   `json:"cancel-job"`
}

// CreateSavepointResponse represents the response body
// used by the create savepoint API
type CreateSavepointResponse struct {
	RequestID string `json:"request-id"`
}

// CreateSavepoint creates a savepoint for a job specified by job ID
func (c FlinkRestClient) CreateSavepoint(jobID string, savepointPath string) (CreateSavepointResponse, error) {
	req := createSavepointRequest{
		TargetDirectory: savepointPath,
		CancelJob:       false,
	}

	reqBody := new(bytes.Buffer)
	json.NewEncoder(reqBody).Encode(req)

	res, err := c.Client.Post(c.constructURL(fmt.Sprintf("jobs/%v/savepoints", jobID)), "application/json", reqBody)
	if err != nil {
		return CreateSavepointResponse{}, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return CreateSavepointResponse{}, err
	}

	if res.StatusCode != 202 {
		return CreateSavepointResponse{}, fmt.Errorf("Unexpected response status %v with body %v", res.StatusCode, string(body[:]))
	}

	response := CreateSavepointResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return CreateSavepointResponse{}, fmt.Errorf("Unable to parse API response as valid JSON: %v", string(body[:]))
	}

	return response, nil
}

// SavepointCreationStatus represents the
// savepoint creation status used by the API
type SavepointCreationStatus struct {
	Id string `json:"id"`
}

// MonitorSavepointCreationResponse represents the response body
// used by the savepoint monitoring API
type MonitorSavepointCreationResponse struct {
	Status SavepointCreationStatus `json:"status"`
}

// MonitorSavepointCreation allows for monitoring the status of a savepoint creation
// identified by the job ID and request ID
func (c FlinkRestClient) MonitorSavepointCreation(jobID string, requestID string) (MonitorSavepointCreationResponse, error) {
	res, err := c.Client.Get(c.constructURL(fmt.Sprintf("jobs/%v/savepoints/%v", jobID, requestID)))
	if err != nil {
		return MonitorSavepointCreationResponse{}, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return MonitorSavepointCreationResponse{}, err
	}

	if res.StatusCode != 200 {
		return MonitorSavepointCreationResponse{}, fmt.Errorf("Unexpected response status %v with body %v", res.StatusCode, string(body[:]))
	}

	response := MonitorSavepointCreationResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return MonitorSavepointCreationResponse{}, fmt.Errorf("Unable to parse API response as valid JSON: %v", string(body[:]))
	}

	return response, nil
}
