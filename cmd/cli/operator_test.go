package main

import (
	"github.com/spf13/afero"
	"ing.com/flink-deployer/cmd/cli/flink"
	"ing.com/flink-deployer/cmd/cli/operations"
)

var mockedDeployError error
var mockedUpdateError error
var mockedRetrieveJobsResponse []flink.Job
var mockedRetrieveJobsError error

type TestOperator struct {
	Filesystem   afero.Fs
	FlinkRestAPI flink.FlinkRestAPI
}

func (t TestOperator) Deploy(d operations.Deploy) error {
	return mockedDeployError
}

func (t TestOperator) Update(u operations.UpdateJob) error {
	return mockedUpdateError
}

func (t TestOperator) RetrieveJobs() ([]flink.Job, error) {
	return mockedRetrieveJobsResponse, mockedRetrieveJobsError
}
