package main

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

/*
 * ListAction
 */
func TestListActionShouldThrowAnErrorWhenTheCommandFails(t *testing.T) {
	mockedExitStatus = -1
	commander = TestCommander{}

	context := cli.Context{}
	err := ListAction(&context)

	assert.EqualError(t, err, "exit status 255")
}

/*
 * DeployAction
 */
func TestDeployActionShouldThrowAnErrorWhenBothTheLocalFilenameAndRemoteFilenameArgumentsAreMissing(t *testing.T) {
	mockedExitStatus = 0
	commander = TestCommander{}

	app := cli.App{}
	set := flag.FlagSet{}
	context := cli.NewContext(&app, &set, nil)
	err := DeployAction(context)

	assert.EqualError(t, err, "both flags 'file-name' and 'remote-file-name' unspecified")
}

func TestDeployActionShouldThrowAnErrorWhenBothTheLocalFilenameAndRemoteFilenameArgumentsAreSet(t *testing.T) {
	mockedExitStatus = 0
	commander = TestCommander{}

	app := cli.App{}
	set := flag.FlagSet{}
	set.String("file-name", "file.jar", "")
	set.String("remote-file-name", "http://www.ing.com/flink-job.jar", "")
	context := cli.NewContext(&app, &set, nil)
	err := DeployAction(context)

	assert.EqualError(t, err, "both flags 'file-name' and 'remote-file-name' specified, only one allowed")
}

func TestDeployActionShouldThrowAnErrorWhenTheCommandFails(t *testing.T) {
	mockedExitStatus = -1
	commander = TestCommander{}

	app := cli.App{}
	set := flag.FlagSet{}
	set.String("file-name", "file.jar", "")
	context := cli.NewContext(&app, &set, nil)
	err := DeployAction(context)

	assert.EqualError(t, err, "an error occurred: exit status 255")
}

/*
 * UpdateAction
 */
func TestUpdateActionShouldThrowAnErrorWhenTheJobnameBaseArgumentIsMissing(t *testing.T) {
	mockedExitStatus = 0
	commander = TestCommander{}

	app := cli.App{}
	set := flag.FlagSet{}
	set.String("file-name", "file.jar", "")
	context := cli.NewContext(&app, &set, nil)
	err := UpdateAction(context)

	assert.EqualError(t, err, "unspecified flag 'job-name-base'")
}

func TestUpdateActionShouldThrowAnErrorWhenBothTheLocalFilenameAndRemoteFilenameArgumentsAreMissing(t *testing.T) {
	mockedExitStatus = 0
	commander = TestCommander{}

	app := cli.App{}
	set := flag.FlagSet{}
	set.String("job-name-base", "Job A", "")
	context := cli.NewContext(&app, &set, nil)
	err := UpdateAction(context)

	assert.EqualError(t, err, "both flags 'file-name' and 'remote-file-name' unspecified")
}

func TestUpdateActionShouldThrowAnErrorWhenBothTheLocalFilenameAndRemoteFilenameArgumentsAreSet(t *testing.T) {
	mockedExitStatus = 0
	commander = TestCommander{}

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
	mockedExitStatus = -1
	commander = TestCommander{}

	app := cli.App{}
	set := flag.FlagSet{}
	set.String("job-name-base", "Job A", "")
	set.String("file-name", "file.jar", "")
	set.String("savepoint-dir", "/savepoints", "")
	context := cli.NewContext(&app, &set, nil)
	err := UpdateAction(context)

	assert.EqualError(t, err, "an error occurred: exit status 255")
}

/*
 * QueryAction
 */
func TestQueryActionShouldThrowAnErrorWhenTheJobnameArgumentIsMissing(t *testing.T) {
	mockedExitStatus = 0
	commander = TestCommander{}

	app := cli.App{}
	set := flag.FlagSet{}
	set.String("file-name", "file.jar", "")
	context := cli.NewContext(&app, &set, nil)
	err := QueryAction(context)

	assert.EqualError(t, err, "unspecified flag 'job-name'")
}

func TestQueryActionShouldThrowAnErrorWhenTheFilenameArgumentIsMissing(t *testing.T) {
	mockedExitStatus = 0
	commander = TestCommander{}

	app := cli.App{}
	set := flag.FlagSet{}
	set.String("job-name", "Job A", "")
	context := cli.NewContext(&app, &set, nil)
	err := QueryAction(context)

	assert.EqualError(t, err, "unspecified flag 'file-name'")
}

func TestQueryActionShouldThrowAnErrorWhenTheMainClassArgumentIsMissing(t *testing.T) {
	mockedExitStatus = 0
	commander = TestCommander{}

	app := cli.App{}
	set := flag.FlagSet{}
	set.String("file-name", "file.jar", "")
	set.String("job-name", "Job A", "")
	context := cli.NewContext(&app, &set, nil)
	err := QueryAction(context)

	assert.EqualError(t, err, "unspecified flag 'main-class'")
}

func TestQueryActionShouldThrowAnErrorWhenTheHighAvailabilityArgumentIsMissing(t *testing.T) {
	mockedExitStatus = 0
	commander = TestCommander{}

	app := cli.App{}
	set := flag.FlagSet{}
	set.String("file-name", "file.jar", "")
	set.String("job-name", "Job A", "")
	set.String("main-class", "com.ing.QueryState", "")
	context := cli.NewContext(&app, &set, nil)
	err := QueryAction(context)

	assert.EqualError(t, err, "unspecified flag 'high-availability'")
}

func TestQueryActionShouldThrowAnErrorWhenHAIsZookeeperAndTheQourumArgumentIsMissing(t *testing.T) {
	mockedExitStatus = 0
	commander = TestCommander{}

	app := cli.App{}
	set := flag.FlagSet{}
	set.String("file-name", "file.jar", "")
	set.String("job-name", "Job A", "")
	set.String("main-class", "com.ing.QueryState", "")
	set.String("high-availability", "zookeeper", "")
	context := cli.NewContext(&app, &set, nil)
	err := QueryAction(context)

	assert.EqualError(t, err, "unspecified flag 'zookeeper-quorum'")
}

func TestQueryActionShouldThrowAnErrorWhenHAIsNoneAndTheJobmanagerAddressIsMissing(t *testing.T) {
	mockedExitStatus = 0
	commander = TestCommander{}

	app := cli.App{}
	set := flag.FlagSet{}
	set.String("file-name", "file.jar", "")
	set.String("job-name", "Job A", "")
	set.String("main-class", "com.ing.QueryState", "")
	set.String("high-availability", "none", "")
	context := cli.NewContext(&app, &set, nil)
	err := QueryAction(context)

	assert.EqualError(t, err, "unspecified flag 'jobmanager-address'")
}

func TestQueryActionShouldThrowAnErrorWhenHAIsNoneAndTheJobmanagerPortIsMissing(t *testing.T) {
	mockedExitStatus = 0
	commander = TestCommander{}

	app := cli.App{}
	set := flag.FlagSet{}
	set.String("file-name", "file.jar", "")
	set.String("job-name", "Job A", "")
	set.String("main-class", "com.ing.QueryState", "")
	set.String("high-availability", "none", "")
	set.String("jobmanager-address", "flink", "")
	context := cli.NewContext(&app, &set, nil)
	err := QueryAction(context)

	assert.EqualError(t, err, "unspecified flag 'jobmanager-port'")
}

func TestQueryActionShouldThrowAnErrorWhenTheCommandFails(t *testing.T) {
	mockedExitStatus = -1
	commander = TestCommander{}

	app := cli.App{}
	set := flag.FlagSet{}
	set.String("file-name", "file.jar", "")
	set.String("job-name", "Job A", "")
	set.String("main-class", "com.ing.QueryState", "")
	set.String("high-availability", "none", "")
	set.String("jobmanager-address", "flink", "")
	set.Int("jobmanager-port", 6123, "")
	context := cli.NewContext(&app, &set, nil)
	err := QueryAction(context)

	assert.EqualError(t, err, "an error occurred: exit status 255")
}
