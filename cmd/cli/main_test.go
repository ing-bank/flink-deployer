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

	assert.EqualError(t, err, "both flags 'filename' and 'remote-filename' unspecified")
}

func TestDeployActionShouldThrowAnErrorWhenBothTheLocalFilenameAndRemoteFilenameArgumentsAreSet(t *testing.T) {
	mockedExitStatus = 0
	commander = TestCommander{}

	app := cli.App{}
	set := flag.FlagSet{}
	set.String("filename", "file.jar", "")
	set.String("remote-filename", "http://www.ing.com/flink-job.jar", "")
	context := cli.NewContext(&app, &set, nil)
	err := DeployAction(context)

	assert.EqualError(t, err, "both flags 'filename' and 'remote-filename' specified, only one allowed")
}

func TestDeployActionShouldThrowAnErrorWhenTheCommandFails(t *testing.T) {
	mockedExitStatus = -1
	commander = TestCommander{}

	app := cli.App{}
	set := flag.FlagSet{}
	set.String("filename", "file.jar", "")
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
	set.String("filename", "file.jar", "")
	context := cli.NewContext(&app, &set, nil)
	err := UpdateAction(context)

	assert.EqualError(t, err, "unspecified flag 'jobname-base'")
}

func TestUpdateActionShouldThrowAnErrorWhenBothTheLocalFilenameAndRemoteFilenameArgumentsAreMissing(t *testing.T) {
	mockedExitStatus = 0
	commander = TestCommander{}

	app := cli.App{}
	set := flag.FlagSet{}
	set.String("jobname-base", "Job A", "")
	context := cli.NewContext(&app, &set, nil)
	err := UpdateAction(context)

	assert.EqualError(t, err, "both flags 'filename' and 'remote-filename' unspecified")
}

func TestUpdateActionShouldThrowAnErrorWhenBothTheLocalFilenameAndRemoteFilenameArgumentsAreSet(t *testing.T) {
	mockedExitStatus = 0
	commander = TestCommander{}

	app := cli.App{}
	set := flag.FlagSet{}
	set.String("jobname-base", "Job A", "")
	set.String("filename", "file.jar", "")
	set.String("remote-filename", "http://www.ing.com/flink-job.jar", "")
	context := cli.NewContext(&app, &set, nil)
	err := UpdateAction(context)

	assert.EqualError(t, err, "both flags 'filename' and 'remote-filename' specified, only one allowed")
}

func TestUpdateActionShouldThrowAnErrorWhenTheCommandFails(t *testing.T) {
	mockedExitStatus = -1
	commander = TestCommander{}

	app := cli.App{}
	set := flag.FlagSet{}
	set.String("jobname-base", "Job A", "")
	set.String("filename", "file.jar", "")
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
	set.String("filename", "file.jar", "")
	context := cli.NewContext(&app, &set, nil)
	err := QueryAction(context)

	assert.EqualError(t, err, "unspecified flag 'jobname'")
}

func TestQueryActionShouldThrowAnErrorWhenTheFilenameArgumentIsMissing(t *testing.T) {
	mockedExitStatus = 0
	commander = TestCommander{}

	app := cli.App{}
	set := flag.FlagSet{}
	set.String("jobname", "Job A", "")
	context := cli.NewContext(&app, &set, nil)
	err := QueryAction(context)

	assert.EqualError(t, err, "unspecified flag 'filename'")
}

func TestQueryActionShouldThrowAnErrorWhenTheMainClassArgumentIsMissing(t *testing.T) {
	mockedExitStatus = 0
	commander = TestCommander{}

	app := cli.App{}
	set := flag.FlagSet{}
	set.String("filename", "file.jar", "")
	set.String("jobname", "Job A", "")
	context := cli.NewContext(&app, &set, nil)
	err := QueryAction(context)

	assert.EqualError(t, err, "unspecified flag 'mainClass'")
}

func TestQueryActionShouldThrowAnErrorWhenTheHighAvailabilityArgumentIsMissing(t *testing.T) {
	mockedExitStatus = 0
	commander = TestCommander{}

	app := cli.App{}
	set := flag.FlagSet{}
	set.String("filename", "file.jar", "")
	set.String("jobname", "Job A", "")
	set.String("mainClass", "com.ing.QueryState", "")
	context := cli.NewContext(&app, &set, nil)
	err := QueryAction(context)

	assert.EqualError(t, err, "unspecified flag 'highAvailability'")
}

func TestQueryActionShouldThrowAnErrorWhenHAIsZookeeperAndTheQourumArgumentIsMissing(t *testing.T) {
	mockedExitStatus = 0
	commander = TestCommander{}

	app := cli.App{}
	set := flag.FlagSet{}
	set.String("filename", "file.jar", "")
	set.String("jobname", "Job A", "")
	set.String("mainClass", "com.ing.QueryState", "")
	set.String("highAvailability", "zookeeper", "")
	context := cli.NewContext(&app, &set, nil)
	err := QueryAction(context)

	assert.EqualError(t, err, "unspecified flag 'zookeeperQuorum'")
}

func TestQueryActionShouldThrowAnErrorWhenHAIsNoneAndTheJobmanagerAddressIsMissing(t *testing.T) {
	mockedExitStatus = 0
	commander = TestCommander{}

	app := cli.App{}
	set := flag.FlagSet{}
	set.String("filename", "file.jar", "")
	set.String("jobname", "Job A", "")
	set.String("mainClass", "com.ing.QueryState", "")
	set.String("highAvailability", "none", "")
	context := cli.NewContext(&app, &set, nil)
	err := QueryAction(context)

	assert.EqualError(t, err, "unspecified flag 'jobmanagerAddress'")
}

func TestQueryActionShouldThrowAnErrorWhenHAIsNoneAndTheJobmanagerPortIsMissing(t *testing.T) {
	mockedExitStatus = 0
	commander = TestCommander{}

	app := cli.App{}
	set := flag.FlagSet{}
	set.String("filename", "file.jar", "")
	set.String("jobname", "Job A", "")
	set.String("mainClass", "com.ing.QueryState", "")
	set.String("highAvailability", "none", "")
	set.String("jobmanagerAddress", "flink", "")
	context := cli.NewContext(&app, &set, nil)
	err := QueryAction(context)

	assert.EqualError(t, err, "unspecified flag 'jobmanagerPort'")
}

func TestQueryActionShouldThrowAnErrorWhenTheCommandFails(t *testing.T) {
	mockedExitStatus = -1
	commander = TestCommander{}

	app := cli.App{}
	set := flag.FlagSet{}
	set.String("filename", "file.jar", "")
	set.String("jobname", "Job A", "")
	set.String("mainClass", "com.ing.QueryState", "")
	set.String("highAvailability", "none", "")
	set.String("jobmanagerAddress", "flink", "")
	set.Int("jobmanagerPort", 6123, "")
	context := cli.NewContext(&app, &set, nil)
	err := QueryAction(context)

	assert.EqualError(t, err, "an error occurred: exit status 255")
}
