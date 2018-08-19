package flink

type FlinkRestAPI interface {
	Cancel(jobID string) error
	CreateSavepoint(jobID string, savepointPath string) (CreateSavepointResponse, error)
	MonitorSavepointCreation(jobID string, requestID string) (MonitorSavepointCreationResponse, error)
	RetrieveJobs() ([]Job, error)
	RunJar(jarID string, entryClass string, jarArgs string, parallelism int, savepointPath string, allowNonRestoredState bool) error
	UploadJar(filename string) (UploadJarResponse, error)
}
