package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestDeployShouldReturnAnErrorIfTheCommandFails(t *testing.T) {
	mockedStdout = ""
	mockedExitStatus = 1
	commander = TestCommander{}

	deploy := Deploy{
		localFilename: "fake.jar",
	}

	err := deploy.execute()

	assert.EqualError(t, err, "exit status 1")
}

func TestDeployShouldAddASavepointArgumentWhenASavepointPathIsSpecified(t *testing.T) {
	mockedStdout = ""
	mockedExitStatus = 0
	commander = TestCommander{}

	deploy := Deploy{
		localFilename: "fake.jar",
		savepointPath: "/savepoint",
	}

	_ = deploy.execute()

	expected := []string{
		"run",
		"-s",
		"/savepoint",
		"fake.jar",
	}
	assert.Equal(t, expected, receivedArguments)
}

func TestDeployShouldAddAnArgumentWhenAllowNonRestorableStateIsTrue(t *testing.T) {
	mockedStdout = ""
	mockedExitStatus = 0
	commander = TestCommander{}

	deploy := Deploy{
		localFilename: "fake.jar",
		allowNonRestorableState: true,
	}

	_, _ = deploy.execute()

	expected := []string{
		"run",
		"-n",
		"fake.jar",
	}
	assert.Equal(t, expected, receivedArguments)
}

func TestDeployShouldReturnAnErrorIfTheFilenameIsNotSpecified(t *testing.T) {
	mockedStdout = ""
	mockedExitStatus = 1
	commander = TestCommander{}

	deploy := Deploy{}

	_, err := deploy.execute()

	assert.EqualError(t, err, "both properties 'remoteFilename' and 'localFilename' are unspecified")
}

func TestDeployShouldAddAllJarArgumentsWhenSupplied(t *testing.T) {
	mockedStdout = ""
	mockedExitStatus = 0
	commander = TestCommander{}

	deploy := Deploy{
		localFilename: "fake.jar",
		jarArgs: "--kafka.bootstrapServers kafka:9092",
	}

	_, _ = deploy.execute()

	expected := []string{
		"run",
		"fake.jar",
		"--kafka.bootstrapServers",
		"kafka:9092",
	}
	assert.Equal(t, expected, receivedArguments)
}

func TestDeployShouldReturnTheCommandOutputWhenTheCommandSucceeds(t *testing.T) {
	mockedStdout = "Success!"
	mockedExitStatus = 0
	commander = TestCommander{}

	deploy := Deploy{
		localFilename: "fake.jar",
	}

	out, err := deploy.execute()

	assert.Equal(t, mockedStdout, string(out))
	assert.Nil(t, err)
}
