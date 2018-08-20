package operations

import "github.com/ing-bank/flink-deployer/cmd/cli/flink"

func (o RealOperator) RetrieveJobs() ([]flink.Job, error) {
	return o.FlinkRestAPI.RetrieveJobs()
}
