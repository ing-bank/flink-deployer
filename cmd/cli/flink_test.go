package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
 * ExtractJobs
 */
func TestExtractJobsShouldReturnNilWhenNoJobsAreRunning(t *testing.T) {
	output := "No jobs running."

	result := ExtractJobs(output)

	assert.Equal(t, 0, len(result))
}

func TestExtractJobsShouldReturnAllJobIdsForDistinctJobNames(t *testing.T) {
	output := `
	------------------ Running/Restarting Jobs -------------------
	15.11.2017 12:23:37 : jobid1 : Job A (RUNNING)
	15.11.2017 12:23:37 : jobid2 : Job B (RUNNING)
	15.11.2017 12:20:37 : jobid3 : Job A (RUNNING)
	--------------------------------------------------------------
	`

	result := ExtractJobs(output)

	assert.Equal(t, 2, len(result))

	val := result["Job A"]

	assert.Equal(t, 2, len(val))
	assert.Equal(t, []string{"jobid1", "jobid3"}, val)
}

/*
 * RetrieveRunningJobIds
 */

func TestRetrieveRunningJobIdsShouldReturnAnErrorWhenTheCommandFails(t *testing.T) {
	mockedStdout = ""
	mockedExitStatus = -1
	commander = TestCommander{}

	jobs, err := RetrieveRunningJobIds("Job A")

	//fmt.Println(err)
	assert.Nil(t, jobs)
	assert.EqualError(t, err, "exit status 255")
}

func TestRetrieveRunningJobIdsShouldReturnAnEmptySliceWhenNoJobsRunning(t *testing.T) {
	mockedStdout = "No running jobs"
	mockedExitStatus = 0
	commander = TestCommander{}

	jobs, err := RetrieveRunningJobIds("Job A")

	assert.Equal(t, 0, len(jobs))
	assert.Nil(t, err)
}

func TestRetrieveRunningJobIdsShouldReturnASliceWithAllJobIdsForTheSpecifiedJobName(t *testing.T) {
	mockedStdout = `
	------------------ Running/Restarting Jobs -------------------
	15.11.2017 12:23:37 : jobid1 : Job A (RUNNING)
	15.11.2017 12:23:37 : jobid2 : Job B (RUNNING)
	15.11.2017 12:20:37 : jobid3 : Job A (RUNNING)
	--------------------------------------------------------------
	`
	mockedExitStatus = 0
	commander = TestCommander{}

	jobs, _ := RetrieveRunningJobIds("Job A")

	assert.Equal(t, 2, len(jobs))
}

func TestRetrieveRunningJobIdsShouldReturnASliceWithAllJobIdsForTheSpecifiedJobNameBase(t *testing.T) {
	mockedStdout = `
	------------------ Running/Restarting Jobs -------------------
	15.11.2017 12:23:37 : jobid1 : Job A (RUNNING)
	15.11.2017 12:23:37 : jobid2 : Job B (RUNNING)
	15.11.2017 12:20:37 : jobid3 : Job A (RUNNING)
	--------------------------------------------------------------
	`
	mockedExitStatus = 0
	commander = TestCommander{}

	jobs, _ := RetrieveRunningJobIds("Job")

	assert.Equal(t, 3, len(jobs))
}

func TestRetrieveRunningJobIdsShouldReturnAnErrorWhenAnUnknownResponseIsReturned(t *testing.T) {
	mockedStdout = "Major error"
	mockedExitStatus = 0
	commander = TestCommander{}

	jobs, err := RetrieveRunningJobIds("Job A")

	assert.Nil(t, jobs)
	assert.EqualError(t, err, "flink list seemed to have failed")
}

/*
 * CancelJobs
 */

func TestCancelJobShouldReturnAnErrorForEmptyJobIds(t *testing.T) {
	mockedStdout = ""
	commander = TestCommander{}
	_, err := CancelJob("")

	assert.EqualError(t, err, "unspecified argument 'jobId'")
}

func TestCancelJobShouldCancelARunningJob(t *testing.T) {
	mockedStdout = "Cancelled!"
	commander = TestCommander{}
	out, _ := CancelJob("jobid1")

	assert.Equal(t, mockedStdout, string(out))
}

/*
 * ListJobs
 */

func TestListJobsShouldReturnAnOverviewOfJobs(t *testing.T) {
	mockedStdout = `
	2017/11/21 10:20:27 Retrieving JobManager.
	Using address flink/172.18.0.6:6123 to connect to JobManager.
	------------------ Running/Restarting Jobs -------------------
	21.11.2017 10:03:58 : 36e33fb85517e6932a2dbed5b82f1836 : ${FLINK_PROGRAM_NAME} (RUNNING)
	--------------------------------------------------------------
	No scheduled jobs.
	`
	commander = TestCommander{}
	out, _ := ListJobs()

	assert.Equal(t, mockedStdout, string(out))
}

/*
 * Savepoint
 */
func TestSavepointShouldReturnAnErrorForEmptyJobIds(t *testing.T) {
	mockedStdout = ""
	commander = TestCommander{}
	_, err := Savepoint("", "/dir/")

	assert.EqualError(t, err, "unspecified argument 'jobId'")
}

func TestSavepointShouldReturnAnErrorForEmptySavepointTargetDir(t *testing.T) {
	mockedStdout = ""
	commander = TestCommander{}
	_, err := Savepoint("jobId", "")

	assert.EqualError(t, err, "unspecified argument 'savepointTargetDir'")
}

func TestSavepointShouldCreateASavepoint(t *testing.T) {
	mockedStdout = "Job saved"
	commander = TestCommander{}
	out, _ := Savepoint("jobid1", "/dir/")

	assert.Equal(t, mockedStdout, string(out))
}
