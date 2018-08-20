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
 * RetrieveLatestSavepoint
 */
func TestRetrieveLatestSavepointShouldReturnAnErrorIfItCannotReadFromDir(t *testing.T) {
	operator := RealOperator{
		Filesystem: afero.NewMemMapFs(),
	}

	files, err := operator.retrieveLatestSavepoint("/savepoints")

	assert.Equal(t, "", files)
	assert.EqualError(t, err, "open /savepoints: file does not exist")
}

func TestRetrieveLatestSavepointShouldReturnAnTheNewestFile(t *testing.T) {
	filesystem := afero.NewMemMapFs()
	filesystem.Mkdir("/savepoints/", 0755)
	afero.WriteFile(filesystem, "/savepoints/savepoint-683b3f-59401d30cfc4", []byte("file a"), 644)
	afero.WriteFile(filesystem, "/savepoints/savepoint-323b3f-59401d30eoe6", []byte("file b"), 644)

	operator := RealOperator{
		Filesystem: filesystem,
	}

	files, err := operator.retrieveLatestSavepoint("/savepoints")

	assert.Equal(t, "/savepoints/savepoint-323b3f-59401d30eoe6", files)
	assert.Nil(t, err)
}

func TestRetrieveLatestSavepointShouldRemoveTheTrailingSlashFromTheSavepointDirectory(t *testing.T) {
	filesystem := afero.NewMemMapFs()
	filesystem.Mkdir("/savepoints/", 0755)
	afero.WriteFile(filesystem, "/savepoints/savepoint-683b3f-59401d30cfc4", []byte("file a"), 644)
	afero.WriteFile(filesystem, "/savepoints/savepoint-323b3f-59401d30eoe6", []byte("file b"), 644)

	operator := RealOperator{
		Filesystem: filesystem,
	}

	files, err := operator.retrieveLatestSavepoint("/savepoints/")

	assert.Equal(t, "/savepoints/savepoint-323b3f-59401d30eoe6", files)
	assert.Nil(t, err)
}

func TestRetrieveLatestSavepointShouldReturnAnErrorWhenDirEmpty(t *testing.T) {
	filesystem := afero.NewMemMapFs()
	filesystem.Mkdir("/savepoints/", 0755)

	operator := RealOperator{
		Filesystem: filesystem,
	}

	files, err := operator.retrieveLatestSavepoint("/savepoints")

	assert.Equal(t, "", files)
	assert.EqualError(t, err, "No savepoints present in directory: /savepoints")
}

/*
 * filterRunningJobs
 */
func TestFilterRunningJobsShouldReturnAnEmptySliceIfNoJobsWereProvided(t *testing.T) {
	operator := RealOperator{}

	res := operator.filterRunningJobs([]flink.Job{})

	assert.Len(t, res, 0)
}

func TestFilterRunningJobsShouldReturnAnEmptySliceIfNoJobsAreRunning(t *testing.T) {
	operator := RealOperator{}

	res := operator.filterRunningJobs([]flink.Job{
		flink.Job{
			Status: "STOPPED",
		},
	})

	assert.Len(t, res, 0)
}

func TestFilterRunningJobsShouldReturnTheRunningJobs(t *testing.T) {
	operator := RealOperator{}

	res := operator.filterRunningJobs([]flink.Job{
		flink.Job{
			Status: "STOPPED",
		},
		flink.Job{
			Status: "RUNNING",
		},
	})

	assert.Len(t, res, 1)
}

/*
 * monitorSavepointCreation
 */
func TestMonitorSavepointCreationShouldReturnAnErrorWhenTheSavepointFailsToBeCreated(t *testing.T) {
	mockedMonitorSavepointCreationError = errors.New("failed")

	operator := RealOperator{
		FlinkRestAPI: TestFlinkRestClient{
			BaseURL: "http://localhost",
			Client:  &http.Client{},
		},
	}

	err := operator.monitorSavepointCreation("job-id", "request-id", 1)

	assert.EqualError(t, err, "failed to create savepoint for job \"job-id\" within 1 seconds")
}

func TestMonitorSavepointCreationShouldReturnNilWhenTheSavepointIsCreated(t *testing.T) {
	mockedMonitorSavepointCreationError = nil
	mockedMonitorSavepointCreationResponse = flink.MonitorSavepointCreationResponse{
		Status: flink.SavepointCreationStatus{
			Id: "COMPLETED",
		},
	}

	operator := RealOperator{
		FlinkRestAPI: TestFlinkRestClient{
			BaseURL: "http://localhost",
			Client:  &http.Client{},
		},
	}

	err := operator.monitorSavepointCreation("job-id", "request-id", 1)

	assert.Nil(t, err)
}

/*
 * UpdateJob
 */
func TestUpdateJobShouldReturnAnErrorWhenTheJobNameBaseIsUndefined(t *testing.T) {
	operator := RealOperator{
		FlinkRestAPI: TestFlinkRestClient{
			BaseURL: "http://localhost",
			Client:  &http.Client{},
		},
	}

	err := operator.Update(UpdateJob{
		LocalFilename: "testdata/sample.jar",
	})

	assert.EqualError(t, err, "unspecified argument 'JobNameBase'")
}

func TestUpdateJobShouldReturnAnErrorWhenTheSavepointDirectoryIsUndefined(t *testing.T) {
	operator := RealOperator{
		FlinkRestAPI: TestFlinkRestClient{
			BaseURL: "http://localhost",
			Client:  &http.Client{},
		},
	}

	err := operator.Update(UpdateJob{
		JobNameBase:   "WordCountStateful",
		LocalFilename: "testdata/sample.jar",
	})

	assert.EqualError(t, err, "unspecified argument 'SavepointDir'")
}

func TestUpdateJobShouldReturnAnErrorWhenRetrievingTheJobsFails(t *testing.T) {
	mockedRetrieveJobsError = errors.New("failed")

	operator := RealOperator{
		FlinkRestAPI: TestFlinkRestClient{
			BaseURL: "http://localhost",
			Client:  &http.Client{},
		},
	}

	err := operator.Update(UpdateJob{
		JobNameBase:   "WordCountStateful",
		LocalFilename: "testdata/sample.jar",
		SavepointDir:  "/data/flink",
	})

	assert.EqualError(t, err, "retrieving jobs failed: failed")
}

func TestUpdateJobShouldReturnAnErrorWhenTheSavepointCannotBeCreated(t *testing.T) {
	mockedRetrieveJobsError = nil
	mockedRetrieveJobsResponse = []flink.Job{
		flink.Job{
			Id:     "Job-A",
			Name:   "WordCountStateful v1.0",
			Status: "RUNNING",
		},
	}
	mockedCreateSavepointError = errors.New("failed")

	operator := RealOperator{
		FlinkRestAPI: TestFlinkRestClient{
			BaseURL: "http://localhost",
			Client:  &http.Client{},
		},
	}

	err := operator.Update(UpdateJob{
		JobNameBase:   "WordCountStateful",
		LocalFilename: "testdata/sample.jar",
		SavepointDir:  "/data/flink",
	})

	assert.EqualError(t, err, "failed to create savepoint for job Job-A due to error: failed")
}

func TestUpdateJobShouldReturnAnErrorWhenTheJobCannotBeCanceled(t *testing.T) {
	mockedRetrieveJobsError = nil
	mockedRetrieveJobsResponse = []flink.Job{
		flink.Job{
			Id:     "Job-A",
			Name:   "WordCountStateful v1.0",
			Status: "RUNNING",
		},
	}
	mockedCreateSavepointError = nil
	mockedCreateSavepointResponse = flink.CreateSavepointResponse{
		RequestID: "request-id",
	}
	mockedMonitorSavepointCreationResponse = flink.MonitorSavepointCreationResponse{
		Status: flink.SavepointCreationStatus{
			Id: "COMPLETED",
		},
	}
	mockedCancelError = errors.New("failed")

	operator := RealOperator{
		FlinkRestAPI: TestFlinkRestClient{
			BaseURL: "http://localhost",
			Client:  &http.Client{},
		},
	}

	err := operator.Update(UpdateJob{
		JobNameBase:   "WordCountStateful",
		LocalFilename: "../testdata/sample.jar",
		SavepointDir:  "/data/flink",
	})

	assert.EqualError(t, err, "job \"Job-A\" failed to cancel due to: failed")
}

func TestUpdateJobShouldReturnNilWhenTheUpdateSucceeds(t *testing.T) {
	filesystem := afero.NewMemMapFs()
	filesystem.Mkdir("/data/flink/", 0755)
	afero.WriteFile(filesystem, "/data/flink/savepoint-683b3f-59401d30cfc4", []byte("file a"), 644)

	mockedRetrieveJobsError = nil
	mockedRetrieveJobsResponse = []flink.Job{
		flink.Job{
			Id:     "Job-A",
			Name:   "WordCountStateful v1.0",
			Status: "RUNNING",
		},
	}
	mockedCreateSavepointError = nil
	mockedCreateSavepointResponse = flink.CreateSavepointResponse{
		RequestID: "request-id",
	}
	mockedMonitorSavepointCreationResponse = flink.MonitorSavepointCreationResponse{
		Status: flink.SavepointCreationStatus{
			Id: "COMPLETED",
		},
	}
	mockedCancelError = nil
	mockedUploadJarResponse = flink.UploadJarResponse{
		Filename: "/data/flink/sample.jar",
		Status:   "success",
	}
	mockedRunJarError = nil

	operator := RealOperator{
		Filesystem: filesystem,
		FlinkRestAPI: TestFlinkRestClient{
			BaseURL: "http://localhost",
			Client:  &http.Client{},
		},
	}

	err := operator.Update(UpdateJob{
		JobNameBase:   "WordCountStateful",
		LocalFilename: "../testdata/sample.jar",
		SavepointDir:  "/data/flink",
	})

	assert.Nil(t, err)
}

func TestUpdateJobShouldReturnAnErrorWhenMultipleRunningJobsAreFound(t *testing.T) {
	mockedRetrieveJobsError = nil
	mockedRetrieveJobsResponse = []flink.Job{
		flink.Job{
			Id:     "Job-A",
			Name:   "WordCountStateful v1.0",
			Status: "RUNNING",
		},
		flink.Job{
			Id:     "Job-B",
			Name:   "WordCountStateful v1.1",
			Status: "RUNNING",
		},
	}

	operator := RealOperator{
		FlinkRestAPI: TestFlinkRestClient{
			BaseURL: "http://localhost",
			Client:  &http.Client{},
		},
	}

	err := operator.Update(UpdateJob{
		JobNameBase:   "WordCountStateful",
		LocalFilename: "testdata/sample.jar",
		SavepointDir:  "/data/flink",
	})

	assert.EqualError(t, err, "job name with base \"WordCountStateful\" has 2 instances running. Aborting update")
}

// func TestUpdateJobShouldExecuteCorrectlyWhenEverythingGoesFine(t *testing.T) {
// 	// our test setup only allows for 1 exec.Comamnd response per test
// 	// so unfortunately we need 1 response which covers all methods
// 	mockedStdout = `
// 	------------------ Running/Restarting Jobs -------------------
// 	15.11.2017 12:23:37 : jobid1 : Job A (RUNNING)
// 	15.11.2017 12:23:37 : jobid2 : Job B (RUNNING)
// 	--------------------------------------------------------------

// 	Retrieving JobManager.
// 		Using address flink-jobmanager/172.26.0.3:6123 to connect to JobManager.
// 		Triggering savepoint for job 683b3f14d75c470de0aaf2b1e83a3158.
// 		Waiting for response...
// 		Savepoint completed. Path: /data/flink/savepoints/savepoint-683b3f-59401d30cfc4
// 		You can resume your program from this savepoint with the run command.
// 	`
// 	mockedExitStatus = 0
// 	commander = TestCommander{}

// 	update := UpdateJob{
// 		JobNameBase:           "Job A",
// 		LocalFilename:         "file.jar",
// 		EntryClass:            "WordCountStateful",
// 		Parallelism:           1,
// 		ProgramArgs:           "--kafka.bootstrapServers kafka:9092",
// 		SavepointDir:          "/data/flink/savepoints",
// 		AllowNonRestoredState: false,
// 	}

// 	err := update.execute()

// 	assert.Nil(t, err)
// }
