package operations

import (
	"errors"
	"log"
	"strings"
)

type Deploy struct {
	RemoteFilename        string
	ApiToken              string
	LocalFilename         string
	EntryClass            string
	Parallelism           int
	ProgramArgs           string
	SavepointPath         string
	AllowNonRestoredState bool
}

func (o RealOperator) extractJarIDFromFilename(filename string) string {
	parts := strings.Split(filename, "/")
	return parts[len(parts)-1]
}

func (o RealOperator) Deploy(d Deploy) error {
	log.Println("Starting deploy")

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
		_, err := downloadFile(d.RemoteFilename, d.ApiToken, filename)
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
