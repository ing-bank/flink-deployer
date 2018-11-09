package operations

import (
	"errors"
	"fmt"
	"log"
	"time"
	"strings"

	"github.com/cenkalti/backoff"
	"github.com/ing-bank/flink-deployer/cmd/cli/flink"
)

// UpdateJob represents the configuration used for
// updating a job on the Flink cluster
type UpdateJob struct {
	JobNameBase           string
	LocalFilename         string
	RemoteFilename        string
	APIToken              string
	EntryClass            string
	Parallelism           int
	ProgramArgs           string
	SavepointDir          string
	AllowNonRestoredState bool
}

func (o RealOperator) filterJobsByRunningAndNameBase(jobs []flink.Job, jobNameBase string) (ret []flink.Job) {
	for _, job := range jobs {
		if job.Status == "RUNNING" && strings.HasPrefix(job.Name, jobNameBase){
			ret = append(ret, job)
		}
	}
	return
}

func (o RealOperator) monitorSavepointCreation(jobID string, requestID string, maxElapsedTime int) (flink.MonitorSavepointCreationResponse, error) {
	tmp := flink.MonitorSavepointCreationResponse{}
	op := func() error {
		log.Println("checking status of savepoint creation")
		res, err := o.FlinkRestAPI.MonitorSavepointCreation(jobID, requestID)
		tmp = res
		if err != nil {
			return err
		}

		switch res.Status.Id {
		case "COMPLETED":
			return nil
		case "IN_PROGRESS":
			return fmt.Errorf("savepoint creation for job \"%v\" is still pending", jobID)
		default:
			return fmt.Errorf("savepoint creation for job \"%v\" returned an unknown status \"%v\"", jobID, res.Status)
		}
	}
	b := &backoff.ExponentialBackOff{
		InitialInterval:     backoff.DefaultInitialInterval,
		RandomizationFactor: backoff.DefaultRandomizationFactor,
		Multiplier:          backoff.DefaultMultiplier,
		MaxInterval:         backoff.DefaultMaxInterval,
		MaxElapsedTime:      time.Duration(maxElapsedTime) * time.Second,
		Clock:               backoff.SystemClock,
	}
	err := backoff.Retry(op, b)
	if err != nil {
		return flink.MonitorSavepointCreationResponse{}, fmt.Errorf("failed to create savepoint for job \"%v\" within %v seconds", jobID, b.MaxElapsedTime.Seconds())
	}

	b.Reset()

	return tmp, nil
}

// Update executes the actual update of a job on the Flink cluster
func (o RealOperator) Update(u UpdateJob) error {
	if len(u.JobNameBase) == 0 {
		return errors.New("unspecified argument 'JobNameBase'")
	}

	log.Printf("starting job update for base name '%v'\n", u.JobNameBase)

	jobs, err := o.FlinkRestAPI.RetrieveJobs()
	if err != nil {
		return fmt.Errorf("retrieving jobs failed: %v", err)
	}

	runningJobs := o.filterJobsByRunningAndNameBase(jobs, u.JobNameBase)

	deploy := Deploy{
		LocalFilename:         u.LocalFilename,
		RemoteFilename:        u.RemoteFilename,
		APIToken:              u.APIToken,
		EntryClass:            u.EntryClass,
		Parallelism:           u.Parallelism,
		ProgramArgs:           u.ProgramArgs,
		AllowNonRestoredState: u.AllowNonRestoredState,
	}
	switch len(runningJobs) {
	case 0:
		return fmt.Errorf("no instance running for job name base \"%v\". Aborting update", u.JobNameBase)
	case 1:
		log.Printf("found exactly 1 running job with base name: \"%v\"", u.JobNameBase)
		job := runningJobs[0]

		log.Printf("creating savepoint for job \"%v\"", job.ID)
		savepointResponse, err := o.FlinkRestAPI.CreateSavepoint(job.ID, u.SavepointDir)
		if err != nil {
			return fmt.Errorf("failed to create savepoint for job %v due to error: %v", job.ID, err)
		}

		savepointStateResponse, err := o.monitorSavepointCreation(job.ID, savepointResponse.RequestID, 60)
		if err != nil {
			return err
		}

		deploy.SavepointPath = savepointStateResponse.Operation.Location

		if len(deploy.SavepointPath) == 0 {
			return errors.New("Savepoint creation failed")
		}

		log.Printf("created savepoint \"%v\"", savepointStateResponse.Operation.Location)
	default:
		return fmt.Errorf("job name with base \"%v\" has %v instances running. Aborting update", u.JobNameBase, len(runningJobs))
	}

	err = o.Deploy(deploy)
	if err != nil {
		return err
	}

	return nil
}
