package operations

import "ing.com/flink-deployer/cmd/cli/flink"

func (o RealOperator) RetrieveJobs() ([]flink.Job, error) {
	return o.FlinkRestAPI.RetrieveJobs()
}
