package main

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/spf13/afero"
)

func RetrieveLatestSavepoint(dir string) (string, error) {
	if strings.HasSuffix(dir, "/") {
		dir = strings.TrimSuffix(dir, "/")
	}

	files, err := afero.ReadDir(filesystem, dir)
	if err != nil {
		return "", err
	}

	var newestFile string
	var newestTime int64 = 0
	for _, f := range files {
		filePath := dir + "/" + f.Name()
		fi, err := filesystem.Stat(filePath)
		if err != nil {
			return "", err
		}
		currTime := fi.ModTime().Unix()
		if currTime > newestTime {
			newestTime = currTime
			newestFile = filePath
		}
	}

	return newestFile, nil
}

func ExtractSavepointPath(output string) (string, error) {
	rgx := regexp.MustCompile("Savepoint completed. Path: file:(.*)\n")
	matches := rgx.FindAllStringSubmatch(output, -1)

	switch len(matches) {
	case 0:
		return "", errors.New("could not extract savepoint path from Flink's output")
	case 1:
		return matches[0][1], nil
	default:
		return "", errors.New("multiple matches for savepoint found")
	}
}

func CreateSavepoint(jobId string, savepointTargetDir string) (string, error) {
	out, err := Savepoint(jobId, savepointTargetDir)
	if err != nil {
		return "", err
	}

	savepoint, err := ExtractSavepointPath(string(out))
	if err != nil {
		return "", err
	}

	if _, err = afero.Exists(filesystem, savepoint); err != nil {
		return "", err
	}

	return savepoint, nil
}

type UpdateJob struct {
	jobNameBase             string
	runArgs                 string
	localFilename           string
	remoteFilename          string
	apiToken                string
	jarArgs                 string
	savepointDirectory      string
	allowNonRestorableState bool
}

func (u UpdateJob) execute() ([]byte, error) {
	if len(u.jobNameBase) == 0 {
		return nil, errors.New("unspecified argument 'jobNameBase'")
	}
	if len(u.savepointDirectory) == 0 {
		return nil, errors.New("unspecified argument 'savepointDirectory'")
	}

	log.Printf("starting job update for base name: %v, and savepoint dir: %v\n", u.jobNameBase, u.savepointDirectory)

	jobIds, err := RetrieveRunningJobIds(u.jobNameBase)
	if err != nil {
		log.Printf("Retrieving the running jobs failed: %v\n", err)
		return nil, err
	}

	deploy := Deploy{
		runArgs:                 u.runArgs,
		localFilename:           u.localFilename,
		remoteFilename:          u.remoteFilename,
		apiToken:                u.apiToken,
		jarArgs:                 u.jarArgs,
		allowNonRestorableState: u.allowNonRestorableState,
	}
	switch len(jobIds) {
	case 0:
		log.Printf("No instance running for job name base \"%v\". Using last available savepoint\n", u.jobNameBase)

		if len(u.savepointDirectory) == 0 {
			return nil, errors.New("cannot retrieve the latest savepoint without specifying the savepoint directory")
		}

		latestSavepoint, err := RetrieveLatestSavepoint(u.savepointDirectory)
		if err != nil {
			log.Printf("Retrieving the latest savepoint failed: %v\n", err)
			return nil, err
		}

		if len(latestSavepoint) != 0 {
			deploy.savepointPath = latestSavepoint
		}
	case 1:
		log.Printf("Found exactly 1 job with base name: %v\n", u.jobNameBase)
		jobId := jobIds[0]

		savepoint, err := CreateSavepoint(jobId, u.savepointDirectory)
		if err != nil {
			return nil, err
		}

		deploy.savepointPath = savepoint

		CancelJob(jobId)
	default:
		return nil, fmt.Errorf("Jobname base \"%v\" has %v instances running", u.jobNameBase, len(jobIds))
	}

	return deploy.execute()
}
