package main

import (
	"github.com/ing-bank/flink-deployer/cmd/cli/flink"
	"github.com/ing-bank/flink-deployer/cmd/cli/operations"
	"github.com/spf13/afero"
)

var mockedDeployError error
var mockedUpdateError error
var mockedTerminateError error
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

func (t TestOperator) Terminate(te operations.TerminateJob) error {
	return mockedUpdateError
}

func (t TestOperator) RetrieveJobs() ([]flink.Job, error) {
	return mockedRetrieveJobsResponse, mockedRetrieveJobsError
}
