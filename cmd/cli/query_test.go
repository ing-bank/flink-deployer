package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryShouldReturnAnErrorWhenRetrievingTheJobsFails(t *testing.T) {
	mockedStdout = ""
	mockedExitStatus = -1
	commander = TestCommander{}

	query := Query{
		jobNameBase: "Job A",
	}

	out, err := query.execute()

	assert.Nil(t, out)
	assert.EqualError(t, err, "exit status 255")
}

func TestQueryShouldReturnAnErrorWhenThereAreNoJobsRunning(t *testing.T) {
	mockedStdout = `
	------------------ Running/Restarting Jobs -------------------
	15.11.2017 12:23:37 : jobid1 : Job B (RUNNING)
	15.11.2017 12:23:37 : jobid2 : Job C (RUNNING)
	--------------------------------------------------------------
	`
	mockedExitStatus = 0
	commander = TestCommander{}

	query := Query{
		jobNameBase: "Job A",
	}

	out, err := query.execute()

	assert.Nil(t, out)
	assert.EqualError(t, err, "Job A is not an active running job base name")
}

func TestQueryShouldReturnTheQueryCommandOutput(t *testing.T) {
	mockedStdout = `
	------------------ Running/Restarting Jobs -------------------
	15.11.2017 12:23:37 : jobid1 : Job A (RUNNING)
	15.11.2017 12:23:37 : jobid2 : Job B (RUNNING)
	--------------------------------------------------------------
	`
	mockedExitStatus = 0
	commander = TestCommander{}

	query := Query{
		jobNameBase:          "Job A",
		filename:             "file.jar",
		mainClass:            "com.ing.QueryState",
		jobManagerRPCAddress: "flink",
		jobManagerRPCPort:    6123,
	}

	out, err := query.execute()

	assert.Equal(t, mockedStdout, string(out))
	assert.Nil(t, err)

	expected := []string{
		"-cp",
		"file.jar",
		"com.ing.QueryState",
		"jobid1",
		"flink",
		"6123",
	}
	assert.Equal(t, expected, receivedArguments)
}

func TestQueryShouldReturnAnErrorWhenThereAreMultipleJobsRunning(t *testing.T) {
	mockedStdout = `
	------------------ Running/Restarting Jobs -------------------
	15.11.2017 12:23:37 : jobid1 : Job A (RUNNING)
	15.11.2017 12:23:37 : jobid2 : Job B (RUNNING)
	15.11.2017 12:20:37 : jobid3 : Job A (RUNNING)
	--------------------------------------------------------------
	`
	mockedExitStatus = 0
	commander = TestCommander{}

	query := Query{
		jobNameBase: "Job A",
	}

	out, err := query.execute()

	assert.Nil(t, out)
	assert.EqualError(t, err, "Job A has 2 instances running")
}
