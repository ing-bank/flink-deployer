package operations

import (
	"errors"
	"net/http"
	"testing"

	"github.com/ing-bank/flink-deployer/cmd/cli/flink"
	"github.com/spf13/afero"

	"github.com/stretchr/testify/assert"
)

/*
 * extractJarIDFromFilename
 */
func TestExtractJarIDFromFilenameShouldReturnThePartAfterTheLastSlash(t *testing.T) {
	operator := RealOperator{}

	jarID := operator.extractJarIDFromFilename("/data/flink/sample.jar")

	assert.Equal(t, "sample.jar", jarID)
}

/*
 * Deploy
 */
func TestDeployShouldReturnAnErrorWhenBothTheSavepointDirAndSavepointPathAreSet(t *testing.T) {
	operator := RealOperator{}

	err := operator.Deploy(Deploy{
		SavepointDir:  "/data/flink",
		SavepointPath: "/data/flink/savepoint-abc",
	})

	assert.EqualError(t, err, "both properties 'SavepointDir' and 'SavepointPath' are specified")
}

func TestDeployShouldReturnAnErrorWhenNeitherTheLocalOrRemoteFileNameAreSet(t *testing.T) {
	operator := RealOperator{}

	err := operator.Deploy(Deploy{})

	assert.EqualError(t, err, "both properties 'RemoteFilename' and 'LocalFilename' are unspecified")
}

func TestDeployShouldReturnAnErrorWhenTheJarUploadFails(t *testing.T) {
	mockedUploadJarError = errors.New("failed")

	operator := RealOperator{
		FlinkRestAPI: TestFlinkRestClient{
			BaseURL: "http://localhost",
			Client:  &http.Client{},
		},
	}

	err := operator.Deploy(Deploy{
		LocalFilename: "testdata/sample.jar",
	})

	assert.EqualError(t, err, "failed")
}

func TestDeployShouldReturnAnErrorWhenTheLatestSavepointCannotBeRetrieved(t *testing.T) {
	filesystem := afero.NewMemMapFs()
	filesystem.Mkdir("/data/flink/", 0755)

	mockedUploadJarResponse = flink.UploadJarResponse{
		Filename: "/data/flink/sample.jar",
		Status:   "success",
	}
	mockedUploadJarError = nil

	operator := RealOperator{
		Filesystem: filesystem,
		FlinkRestAPI: TestFlinkRestClient{
			BaseURL: "http://localhost",
			Client:  &http.Client{},
		},
	}

	err := operator.Deploy(Deploy{
		LocalFilename: "testdata/sample.jar",
		SavepointDir:  "/data/flink",
	})

	assert.EqualError(t, err, "retrieving the latest savepoint failed: No savepoints present in directory: /data/flink")
}

func TestDeployShouldReturnAnErrorWhenTheJarRunFails(t *testing.T) {
	mockedUploadJarResponse = flink.UploadJarResponse{
		Filename: "/data/flink/sample.jar",
		Status:   "success",
	}
	mockedUploadJarError = nil
	mockedRunJarError = errors.New("failed")

	operator := RealOperator{
		FlinkRestAPI: TestFlinkRestClient{
			BaseURL: "http://localhost",
			Client:  &http.Client{},
		},
	}

	err := operator.Deploy(Deploy{
		LocalFilename: "testdata/sample.jar",
	})

	assert.EqualError(t, err, "failed")
}

func TestDeployShouldReturnNilWhenTheDeploySucceeds(t *testing.T) {
	mockedUploadJarResponse = flink.UploadJarResponse{
		Filename: "/data/flink/sample.jar",
		Status:   "success",
	}
	mockedUploadJarError = nil
	mockedRunJarError = nil

	operator := RealOperator{
		FlinkRestAPI: TestFlinkRestClient{
			BaseURL: "http://localhost",
			Client:  &http.Client{},
		},
	}

	err := operator.Deploy(Deploy{
		LocalFilename: "testdata/sample.jar",
	})

	assert.Nil(t, err)
}
