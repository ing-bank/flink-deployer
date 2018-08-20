package operations

import (
	"github.com/ing-bank/flink-deployer/cmd/cli/flink"
	"github.com/spf13/afero"
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
