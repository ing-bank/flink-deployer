package main

import (
	"errors"
	"flag"
	"os"
	"testing"

	"github.com/ing-bank/flink-deployer/cmd/cli/flink"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

/*
 * Get API Timeout
 */
func TestGetAPITimeOutSecondsShouldReturnTheDefaultValueIfTheEnvVarIsUnset(t *testing.T) {
	timeout, _ := getAPITimeoutSeconds()
	assert.Equal(t, int64(10), timeout)
}

func TestGetAPITimeOutSecondsShouldReturnTheParsedValueFromTheEnvVar(t *testing.T) {
	os.Setenv("FLINK_API_TIMEOUT_SECONDS", "20")
	timeout, _ := getAPITimeoutSeconds()
	assert.Equal(t, int64(20), timeout)
}

func TestGetAPITimeOutSecondsShouldReturnAnErrorIfTheEnvVarValueCannotBeParsed(t *testing.T) {
	os.Setenv("FLINK_API_TIMEOUT_SECONDS", "bla")
	_, err := getAPITimeoutSeconds()
	assert.EqualError(t, err, "strconv.ParseInt: parsing \"bla\": invalid syntax")
}

/*
 * ListAction
 */
func TestListActionShouldReturnAnErrorWhenTheAPIFails(t *testing.T) {
	mockedRetrieveJobsError = errors.New("failed")

	operator = TestOperator{}

	context := cli.Context{}
	err := ListAction(&context)

	assert.EqualError(t, err, "failed to list jobs: failed")
}

func TestListActionShouldReturnNilWhenTheAPISucceeds(t *testing.T) {
	mockedRetrieveJobsResponse = []flink.Job{
		flink.Job{
			ID:     "1",
			Name:   "Job A",
			Status: "RUNNING",
		},
	}
	mockedRetrieveJobsError = nil

	operator = TestOperator{}

	context := cli.Context{}
	err := ListAction(&context)

	assert.Nil(t, err)
}

/*
 * TerminateAction
 */
func TestTerminateActionShouldThrowAnErrorWhenJobNameBaseMissing(t *testing.T) {
	mockedUpdateError = nil
	operator = TestOperator{}

	app := cli.App{}
	set := flag.FlagSet{}
	context := cli.NewContext(&app, &set, nil)
	err := TerminateAction(context)

	assert.EqualError(t, err, "unspecified flag 'job-name-base'")
}

func TestTerminateActionShouldThrowAnErrorWhenModeUnknown(t *testing.T) {
	operator = TestOperator{}

	app := cli.App{}
	set := flag.FlagSet{}
	set.String("job-name-base", "a job", "")
	set.String("mode", "not_exist_mode", "")
	context := cli.NewContext(&app, &set, nil)
	err := TerminateAction(context)

	assert.EqualError(t, err, "unknown flag 'mode', only `cancel` and `stop` supported")
}

/*
 * DeployAction
 */
func TestDeployActionShouldThrowAnErrorWhenBothTheLocalFilenameAndRemoteFilenameArgumentsAreMissing(t *testing.T) {
	operator = TestOperator{}

	app := cli.App{}
	set := flag.FlagSet{}
	context := cli.NewContext(&app, &set, nil)
	err := DeployAction(context)

	assert.EqualError(t, err, "both flags 'file-name' and 'remote-file-name' unspecified")
}

func TestDeployActionShouldThrowAnErrorWhenBothTheLocalFilenameAndRemoteFilenameArgumentsAreSet(t *testing.T) {
	operator = TestOperator{}

	app := cli.App{}
	set := flag.FlagSet{}
	set.String("file-name", "file.jar", "")
	set.String("remote-file-name", "http://www.ing.com/flink-job.jar", "")
	context := cli.NewContext(&app, &set, nil)
	err := DeployAction(context)

	assert.EqualError(t, err, "both flags 'file-name' and 'remote-file-name' specified, only one allowed")
}

func TestDeployActionShouldThrowAnErrorWhenBothTheSavepointDirAndSavepointPathArgumentsAreSet(t *testing.T) {
	operator = TestOperator{}

	app := cli.App{}
	set := flag.FlagSet{}
	set.String("file-name", "file.jar", "")
	set.String("savepoint-dir", "/data/flink", "")
	set.String("savepoint-path", "/data/flink/savepoint-abc", "")
	context := cli.NewContext(&app, &set, nil)
	err := DeployAction(context)

	assert.EqualError(t, err, "both flags 'savepoint-dir' and 'savepoint-path' specified, only one allowed")
}

func TestDeployActionShouldThrowAnErrorWhenTheCommandFails(t *testing.T) {
	mockedDeployError = errors.New("failed")
	operator = TestOperator{}

	app := cli.App{}
	set := flag.FlagSet{}
	set.String("file-name", "file.jar", "")
	context := cli.NewContext(&app, &set, nil)
	err := DeployAction(context)

	assert.EqualError(t, err, "an error occurred: failed")
}

/*
 * UpdateAction
 */
func TestUpdateActionShouldThrowAnErrorWhenTheJobnameBaseArgumentIsMissing(t *testing.T) {
	mockedUpdateError = nil
	operator = TestOperator{}

	app := cli.App{}
	set := flag.FlagSet{}
	set.String("file-name", "file.jar", "")
	context := cli.NewContext(&app, &set, nil)
	err := UpdateAction(context)

	assert.EqualError(t, err, "unspecified flag 'job-name-base'")
}

func TestUpdateActionShouldThrowAnErrorWhenBothTheLocalFilenameAndRemoteFilenameArgumentsAreMissing(t *testing.T) {
	mockedUpdateError = nil
	operator = TestOperator{}

	app := cli.App{}
	set := flag.FlagSet{}
	set.String("job-name-base", "Job A", "")
	context := cli.NewContext(&app, &set, nil)
	err := UpdateAction(context)

	assert.EqualError(t, err, "both flags 'file-name' and 'remote-file-name' unspecified")
}

func TestUpdateActionShouldThrowAnErrorWhenBothTheLocalFilenameAndRemoteFilenameArgumentsAreSet(t *testing.T) {
	mockedUpdateError = nil
	operator = TestOperator{}

	app := cli.App{}
	set := flag.FlagSet{}
	set.String("job-name-base", "Job A", "")
	set.String("file-name", "file.jar", "")
	set.String("remote-file-name", "http://www.ing.com/flink-job.jar", "")
	context := cli.NewContext(&app, &set, nil)
	err := UpdateAction(context)

	assert.EqualError(t, err, "both flags 'file-name' and 'remote-file-name' specified, only one allowed")
}

func TestUpdateActionShouldThrowAnErrorWhenTheCommandFails(t *testing.T) {
	mockedUpdateError = errors.New("failed")
	operator = TestOperator{}

	app := cli.App{}
	set := flag.FlagSet{}
	set.String("job-name-base", "Job A", "")
	set.String("file-name", "file.jar", "")
	set.String("savepoint-dir", "/savepoints", "")
	context := cli.NewContext(&app, &set, nil)
	err := UpdateAction(context)

	assert.EqualError(t, err, "an error occurred: failed")
}
