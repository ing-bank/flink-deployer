package main

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

type Deploy struct {
	runArgs                 string
	remoteFilename          string
	apiToken                string
	localFilename           string
	jarArgs                 string
	savepointPath           string
	allowNonRestorableState bool
}

func (d Deploy) execute() ([]byte, error) {
	log.Println("Starting deploy")

	args := []string{
		"run",
	}

	if len(d.savepointPath) > 0 {
		log.Printf("Using savepoint for deployment: %v", d.savepointPath)
		args = append(args, []string{"-s", d.savepointPath}...)
	}

	if d.allowNonRestorableState == true {
		log.Printf("Allowing non restorable state")
		args = append(args, "-n")
	}

	if len(d.runArgs) != 0 {
		runArgs := strings.Split(d.runArgs, " ")

		for _, v := range runArgs {
			if len(v) == 0 {
				continue
			}
			args = append(args, fmt.Sprintf("%v", v))
		}
	}

	if len(d.remoteFilename) == 0 && len(d.localFilename) == 0 {
		return nil, errors.New("both properties 'remoteFilename' and 'localFilename' are unspecified")
	}

	if len(d.remoteFilename) > 0 {
		filename := "/tmp/job.jar"
		_, err := downloadFile(d.remoteFilename, d.apiToken, filename)
		if err != nil {
			return nil, err
		}
		args = append(args, filename)
	}

	if len(d.localFilename) > 0 {
		args = append(args, d.localFilename)
	}

	if len(d.jarArgs) != 0 {
		jarArgs := strings.Split(d.jarArgs, " ")

		for _, v := range jarArgs {
			if len(v) == 0 {
				continue
			}
			args = append(args, fmt.Sprintf("%v", v))
		}
	}

	log.Println("Deploying job")
	log.Printf("Arguments: %v\n", args)

	out, err := commander.CombinedOutput("flink", args...)
	if err != nil {
		return out, err
	}

	return out, nil
}
