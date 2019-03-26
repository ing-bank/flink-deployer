package operations

import (
	"errors"
	"fmt"
)

// UpdateJob represents the configuration used for
// updating a job on the Flink cluster
type TerminateJob struct {
	JobNameBase string
	Mode        string
}

// Update executes the actual update of a job on the Flink cluster
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
