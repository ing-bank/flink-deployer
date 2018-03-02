package main

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

/*
 * RetrieveLatestSavepoint
 */
func TestRetrieveLatestSavepointShouldReturnAnErrorIfItCannotReadFromDir(t *testing.T) {
	filesystem = afero.NewMemMapFs()

	files, err := RetrieveLatestSavepoint("/savepoints")

	assert.Equal(t, "", files)
	assert.EqualError(t, err, "open /savepoints: file does not exist")
}

func TestRetrieveLatestSavepointShouldReturnAnTheNewestFile(t *testing.T) {
	filesystem = afero.NewMemMapFs()
	filesystem.Mkdir("/savepoints/", 0755)
	afero.WriteFile(filesystem, "/savepoints/savepoint-683b3f-59401d30cfc4", []byte("file a"), 644)
	afero.WriteFile(filesystem, "/savepoints/savepoint-323b3f-59401d30eoe6", []byte("file b"), 644)

	files, err := RetrieveLatestSavepoint("/savepoints")

	assert.Equal(t, "/savepoints/savepoint-323b3f-59401d30eoe6", files)
	assert.Nil(t, err)
}

func TestRetrieveLatestSavepointShouldRemoveTheTrailingSlashFromTheSavepointDirectory(t *testing.T) {
	filesystem = afero.NewMemMapFs()
	filesystem.Mkdir("/savepoints/", 0755)
	afero.WriteFile(filesystem, "/savepoints/savepoint-683b3f-59401d30cfc4", []byte("file a"), 644)
	afero.WriteFile(filesystem, "/savepoints/savepoint-323b3f-59401d30eoe6", []byte("file b"), 644)

	files, err := RetrieveLatestSavepoint("/savepoints/")

	assert.Equal(t, "/savepoints/savepoint-323b3f-59401d30eoe6", files)
	assert.Nil(t, err)
}

func TestRetrieveLatestSavepointShouldReturnAnErrorWhenDirEmpty(t *testing.T) {
	filesystem = afero.NewMemMapFs()
	filesystem.Mkdir("/savepoints/", 0755)

	files, err := RetrieveLatestSavepoint("/savepoints")

	assert.Equal(t, "", files)
	assert.EqualError(t, err, "No savepoints present in directory: /savepoints")
}

/*
 * ExtractSavepointPath
 */
func TestExtractSavepointPathShouldExtractPath(t *testing.T) {
	out, _ := ExtractSavepointPath(`
		Retrieving JobManager.
		Using address flink-jobmanager/172.26.0.3:6123 to connect to JobManager.
		Triggering savepoint for job 683b3f14d75c470de0aaf2b1e83a3158.
		Waiting for response...
		Savepoint completed. Path: file:/data/flink/savepoints/savepoint-683b3f-59401d30cfc4
		You can resume your program from this savepoint with the run command.
		`)

	assert.Equal(t, "/data/flink/savepoints/savepoint-683b3f-59401d30cfc4", out)
}

func TestExtractSavepointPathShouldReturnAnErrorIfItCannotExtractThePath(t *testing.T) {

	out, err := ExtractSavepointPath(`
		Error: The outue for the Job ID is not a outid ID.

		Use the help option (-h or --help) to get help on the command.
		`)

	assert.Equal(t, "", out)
	assert.EqualError(t, err, "could not extract savepoint path from Flink's output")
}

func TestExtractSavepointPathShouldReturnAnErrorIfMultiplePathsAreExtracted(t *testing.T) {
	out, err := ExtractSavepointPath(`
		Retrieving JobManager.
		Using address flink-jobmanager/172.26.0.3:6123 to connect to JobManager.
		Triggering savepoint for job 683b3f14d75c470de0aaf2b1e83a3158.
		Waiting for response...
		Savepoint completed. Path: file:/data/flink/savepoints/savepoint-683b3f-59401d30cfc4
		Savepoint completed. Path: file:/data/flink/savepoints/savepoint-883b3f-59401d30cfc1
		You can resume your program from this savepoint with the run command.
		`)

	assert.Equal(t, "", out)
	assert.EqualError(t, err, "multiple matches for savepoint found")
}

/*
 * CreateSavepoint
 */

func TestCreateSavepointShouldReturnAnErrorWhenCreatingTheSavepointFails(t *testing.T) {
	mockedExitStatus = -1
	commander = TestCommander{}

	out, err := CreateSavepoint("182b71aebf67191683b6917ce95a1f34", "/dir/")

	assert.Equal(t, "", out)
	assert.EqualError(t, err, "exit status 255")
}

func TestCreateSavepointShouldReturnAnErrorWhenExtractingTheSavepointPathFails(t *testing.T) {
	mockedStdout = `
		Retrieving JobManager.
		Using address flink-jobmanager/172.26.0.3:6123 to connect to JobManager.
		Triggering savepoint for job 683b3f14d75c470de0aaf2b1e83a3158.
		Waiting for response...
		Savepoint completed. Path: file:/data/flink/savepoints/savepoint-683b3f-59401d30cfc4
		Savepoint completed. Path: file:/data/flink/savepoints/savepoint-883b3f-59401d30cfc1
		You can resume your program from this savepoint with the run command.
	`
	mockedExitStatus = 0
	commander = TestCommander{}

	out, err := CreateSavepoint("182b71aebf67191683b6917ce95a1f34", "/data/flink/savepoints")

	assert.Equal(t, "", out)
	assert.EqualError(t, err, "multiple matches for savepoint found")
}

func TestCreateSavepointShouldReturnTheSavepointPathIfAllGoesWell(t *testing.T) {
	mockedStdout = `
		Retrieving JobManager.
		Using address flink-jobmanager/172.26.0.3:6123 to connect to JobManager.
		Triggering savepoint for job 683b3f14d75c470de0aaf2b1e83a3158.
		Waiting for response...
		Savepoint completed. Path: file:/data/flink/savepoints/savepoint-683b3f-59401d30cfc4
		You can resume your program from this savepoint with the run command.
	`
	mockedExitStatus = 0
	commander = TestCommander{}
	filesystem = afero.NewMemMapFs()

	out, err := CreateSavepoint("182b71aebf67191683b6917ce95a1f34", "/data/flink/savepoints")

	assert.Equal(t, "/data/flink/savepoints/savepoint-683b3f-59401d30cfc4", out)
	assert.Nil(t, err)
}

/*
 * UpdateJob
 */

func TestUpdateJobShouldReturnAnErrorWhenTheJobNameBaseIsUndefined(t *testing.T) {
	mockedStdout = `No running jobs`
	mockedExitStatus = 0
	commander = TestCommander{}

	update := UpdateJob{}

	out, err := update.execute()

	assert.Nil(t, out)
	assert.EqualError(t, err, "unspecified argument 'jobNameBase'")
}

func TestUpdateJobShouldReturnAnErrorWhenTheSavepointDirectoryIsUndefined(t *testing.T) {
	mockedStdout = `No running jobs`
	mockedExitStatus = 0
	commander = TestCommander{}

	update := UpdateJob{
		jobNameBase: "Job A",
	}

	out, err := update.execute()

	assert.Nil(t, out)
	assert.EqualError(t, err, "unspecified argument 'savepointDirectory'")
}

func TestUpdateJobShouldExecuteCorrectlyWhenEverythingGoesFine(t *testing.T) {
	// our test setup only allows for 1 exec.Comamnd response per test
	// so unfortunately we need 1 response which covers all methods
	mockedStdout = `
	------------------ Running/Restarting Jobs -------------------
	15.11.2017 12:23:37 : jobid1 : Job A (RUNNING)
	15.11.2017 12:23:37 : jobid2 : Job B (RUNNING)
	--------------------------------------------------------------

	Retrieving JobManager.
		Using address flink-jobmanager/172.26.0.3:6123 to connect to JobManager.
		Triggering savepoint for job 683b3f14d75c470de0aaf2b1e83a3158.
		Waiting for response...
		Savepoint completed. Path: file:/data/flink/savepoints/savepoint-683b3f-59401d30cfc4
		You can resume your program from this savepoint with the run command.
	`
	mockedExitStatus = 0
	commander = TestCommander{}

	update := UpdateJob{
		jobNameBase:             "Job A",
		runArgs:                 "-p 1 -d",
		localFilename:           "file.jar",
		jarArgs:                 "--kafka.bootstrapServers kafka:9092",
		savepointDirectory:      "/data/flink/savepoints",
		allowNonRestorableState: false,
	}

	out, err := update.execute()

	assert.Equal(t, mockedStdout, string(out))
	assert.Nil(t, err)
}

func TestUpdateJobShouldReturnAnErrorWhenMultipleRunningJobsAreFound(t *testing.T) {
	// our test setup only allows for 1 exec.Comamnd response per test
	// so unfortunately we need 1 response which covers all methods
	mockedStdout = `
	------------------ Running/Restarting Jobs -------------------
	15.11.2017 12:23:37 : jobid1 : Job A (RUNNING)
	15.11.2017 12:23:37 : jobid2 : Job B (RUNNING)
	15.11.2017 12:20:37 : jobid3 : Job A (RUNNING)
	--------------------------------------------------------------

	Retrieving JobManager.
		Using address flink-jobmanager/172.26.0.3:6123 to connect to JobManager.
		Triggering savepoint for job 683b3f14d75c470de0aaf2b1e83a3158.
		Waiting for response...
		Savepoint completed. Path: file:/data/flink/savepoints/savepoint-683b3f-59401d30cfc4
		You can resume your program from this savepoint with the run command.
	`
	mockedExitStatus = 0
	commander = TestCommander{}

	update := UpdateJob{
		jobNameBase:             "Job A",
		runArgs:                 "-p 1 -d",
		localFilename:           "file.jar",
		jarArgs:                 "--kafka.bootstrapServers kafka:9092",
		savepointDirectory:      "/data/flink/savepoints",
		allowNonRestorableState: false,
	}

	out, err := update.execute()

	assert.Nil(t, out)
	assert.EqualError(t, err, "Jobname base \"Job A\" has 2 instances running")
}
