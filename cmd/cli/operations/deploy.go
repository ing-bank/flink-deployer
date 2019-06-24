package operations

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

// Deploy represents the configuration used for
// deploying a job to the Flink cluster
type Deploy struct {
	RemoteFilename        string
	APIToken              string
	LocalFilename         string
	EntryClass            string
	Parallelism           int
	ProgramArgs           []string
	SavepointDir          string
	SavepointPath         string
	AllowNonRestoredState bool
}

func (o RealOperator) extractJarIDFromFilename(filename string) string {
	parts := strings.Split(filename, "/")
	return parts[len(parts)-1]
}

// Deploy executes the actual deployment to the Flink cluster
func (o RealOperator) Deploy(d Deploy) error {
	log.Println("Starting deploy")

	if len(d.SavepointDir) > 0 && len(d.SavepointPath) > 0 {
		return errors.New("both properties 'SavepointDir' and 'SavepointPath' are specified")
	}

	if len(d.SavepointDir) > 0 {
		log.Printf("Using savepoint directory to retrieve the latest savepoint: %v", d.SavepointDir)

		latestSavepoint, err := o.retrieveLatestSavepoint(d.SavepointDir)
		if err != nil {
			return fmt.Errorf("retrieving the latest savepoint failed: %v", err)
		}

		if len(latestSavepoint) != 0 {
			d.SavepointPath = latestSavepoint
		}
	}

	if len(d.SavepointPath) > 0 {
		log.Printf("Using savepoint for deployment: %v", d.SavepointPath)
	}

	if d.AllowNonRestoredState == true {
		log.Printf("Allowing non restorable state")
	}

	if len(d.RemoteFilename) == 0 && len(d.LocalFilename) == 0 {
		return errors.New("both properties 'RemoteFilename' and 'LocalFilename' are unspecified")
	}

	var filename string

	if len(d.RemoteFilename) > 0 {
		filename = "/tmp/job.jar"
		_, err := downloadFile(d.RemoteFilename, d.APIToken, filename)
		if err != nil {
			return err
		}
	}

	if len(d.LocalFilename) > 0 {
		filename = d.LocalFilename
	}

	log.Println("Uploading JAR file")
	uploadResponse, err := o.FlinkRestAPI.UploadJar(filename)
	if err != nil {
		return err
	}

	jarID := o.extractJarIDFromFilename(uploadResponse.Filename)

	log.Println("Running job")
	err = o.FlinkRestAPI.RunJar(jarID, d.EntryClass, d.ProgramArgs, d.Parallelism, d.SavepointPath, d.AllowNonRestoredState)
	if err != nil {
		return err
	}

	return nil
}
