package operations

import (
	"github.com/ing-bank/flink-deployer/cmd/cli/flink"
	"github.com/spf13/afero"
)

// Operator is an interface which contains all the functionality
// that the deployer exposes
type Operator interface {
	Deploy(d Deploy) error
	Update(u UpdateJob) error
	RetrieveJobs() ([]flink.Job, error)
	Terminate(t TerminateJob) error
}

// RealOperator is the Operator used in the production code
type RealOperator struct {
	Filesystem   afero.Fs
	FlinkRestAPI flink.FlinkRestAPI
}
