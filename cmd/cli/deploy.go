package main

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

type Deploy struct {
	remoteFilename          string
	apiToken                string
	localFilename           string
	entryClass              string
	parallelism             int
	jarArgs                 string
	savepointPath           string
	allowNonRestorableState bool
}

func findJarIDForFilename(filename string, jars []Jar) (string, error) {
	for _, jar := range jars {
		log.Printf("Jar found: %v\n", jar.name)
		if jar.name == filename {
			return jar.id, nil
		}
	}

	return "", fmt.Errorf("Unable to find JAR with the filename: %v", filename)
}

func (d Deploy) execute() error {
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

	if len(d.remoteFilename) == 0 && len(d.localFilename) == 0 {
		return errors.New("both properties 'remoteFilename' and 'localFilename' are unspecified")
	}

	var filename string

	if len(d.remoteFilename) > 0 {
		filename = "/tmp/job.jar"
		_, err := downloadFile(d.remoteFilename, d.apiToken, filename)
		if err != nil {
			return err
		}
	}

	if len(d.localFilename) > 0 {
		filename = d.localFilename
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

	uploadResponse, err := flinkRestClient.uploadJar(filename)
	if err != nil {
		return err
	}

	jars, err := flinkRestClient.retrieveJars()
	if err != nil {
		return err
	}

	jarID, err := findJarIDForFilename(uploadResponse.filename, jars)
	if err != nil {
		return err
	}

	err = flinkRestClient.runJar(jarID, d.jarArgs, d.savepointPath, d.allowNonRestorableState)
	if err != nil {
		return err
	}

	return nil
}
