package operations

import (
	"net/http"

	"github.com/ing-bank/flink-deployer/cmd/cli/flink"
)

/*
 * Flink REST API mocking
 */

var mockedTerminateError error
var mockedCreateSavepointResponse flink.CreateSavepointResponse
var mockedCreateSavepointError error
var mockedMonitorSavepointCreationResponse flink.MonitorSavepointCreationResponse
var mockedMonitorSavepointCreationError error
var mockedRetrieveJobsResponse []flink.Job
var mockedRetrieveJobsError error
var mockedRunJarError error
var mockedUploadJarResponse flink.UploadJarResponse
var mockedUploadJarError error

type TestFlinkRestClient struct {
	BaseURL string
	Client  *http.Client
}

func (c TestFlinkRestClient) Terminate(jobID string, mode string) error {
	return mockedTerminateError
}
func (c TestFlinkRestClient) CreateSavepoint(jobID string, savepointPath string) (flink.CreateSavepointResponse, error) {
	return mockedCreateSavepointResponse, mockedCreateSavepointError
}
func (c TestFlinkRestClient) MonitorSavepointCreation(jobID string, requestID string) (flink.MonitorSavepointCreationResponse, error) {
	return mockedMonitorSavepointCreationResponse, mockedMonitorSavepointCreationError
}
func (c TestFlinkRestClient) RetrieveJobs() ([]flink.Job, error) {
	return mockedRetrieveJobsResponse, mockedRetrieveJobsError
}
func (c TestFlinkRestClient) RunJar(jarID string, entryClass string, jarArgs []string, parallelism int, savepointPath string, allowNonRestoredState bool) error {
	return mockedRunJarError
}
func (c TestFlinkRestClient) UploadJar(filename string) (flink.UploadJarResponse, error) {
	return mockedUploadJarResponse, mockedUploadJarError
}

func constructTestClient() flink.FlinkRestAPI {
	return TestFlinkRestClient{
		BaseURL: "http://localhost",
		Client:  &http.Client{},
	}
}
