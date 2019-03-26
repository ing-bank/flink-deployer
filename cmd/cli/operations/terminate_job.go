package operations

import (
	"errors"
	"fmt"
)

// TerminateJob represents the configuration used for
// terminate a job on the Flink cluster
type TerminateJob struct {
	JobNameBase string
	Mode        string
}

// Terminate executes the actual termination of a job on the Flink cluster
func (o RealOperator) Terminate(t TerminateJob) error {
	if len(t.JobNameBase) == 0 {
		return errors.New("unspecified argument 'JobNameBase'")
	}

	err := o.FlinkRestAPI.Terminate(t.JobNameBase, t.Mode)
	if err != nil {
		return fmt.Errorf("job \"%v\" failed to terminate due to: %v", t.JobNameBase, err)
	}

	return nil
}
