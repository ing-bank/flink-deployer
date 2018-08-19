package operations

import (
	"github.com/spf13/afero"
	"ing.com/flink-deployer/cmd/cli/flink"
)

type Operator interface {
	Deploy(d Deploy) error
	Update(u UpdateJob) error
	RetrieveJobs() ([]flink.Job, error)
}

type RealOperator struct {
	Filesystem   afero.Fs
	FlinkRestAPI flink.FlinkRestAPI
}
