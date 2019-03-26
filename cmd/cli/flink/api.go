package flink

// FlinkRestAPI is an interface representing the ability to execute
// multiple HTTP requests against the Apache Flink API.
type FlinkRestAPI interface {
	Cancel(jobID string) error
	Terminate(jobID string, mode string) error
	CreateSavepoint(jobID string, savepointPath string) (CreateSavepointResponse, error)
	MonitorSavepointCreation(jobID string, requestID string) (MonitorSavepointCreationResponse, error)
	RetrieveJobs() ([]Job, error)
	RunJar(jarID string, entryClass string, jarArgs string, parallelism int, savepointPath string, allowNonRestoredState bool) error
	UploadJar(filename string) (UploadJarResponse, error)
}
