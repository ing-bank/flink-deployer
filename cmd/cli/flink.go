package main

import (
	"errors"
	"log"
	"regexp"
	"strings"
)

func ExtractJobs(output string) map[string][]string {
	jobNameRgx := regexp.MustCompile("[[0-9:.]+ [[0-9:.]+ : ([0-9a-z]+) : (.*) \\(([A-Z]+)\\)\\n")
	jobNameMatches := jobNameRgx.FindAllStringSubmatch(output, -1)

	jobs := make(map[string][]string)
	for _, v := range jobNameMatches {
		jobId := v[1]
		jobName := v[2]
		if val, ok := jobs[jobName]; ok {
			jobs[jobName] = append(val, jobId)
		} else {
			jobs[jobName] = []string{jobId}
		}
	}

	return jobs
}

func RetrieveRunningJobIds(jobNameBase string) ([]string, error) {
	out, err := commander.CombinedOutput("flink", "list", "-r")
	if err != nil {
		return nil, err
	}

	output := string(out)
	if strings.Contains(output, "No running jobs") {
		log.Printf("No running job found for base name %v. Continuing with deploy\n", jobNameBase)
		return nil, nil
	} else if strings.Contains(output, "Running/Restarting Jobs") {
		jobs := ExtractJobs(output)

		// Extract all entries in jobs map with jobname that is prepended by jobNameBase
		var jobIds []string
		for jobName, jobIDList := range jobs {
			if strings.HasPrefix(jobName, jobNameBase) {
				jobIds = append(jobIds, jobIDList...)
			}
		}
		return jobIds, nil
	} else {
		return nil, errors.New("flink list seemed to have failed")
	}
}

func CancelJob(jobId string) ([]byte, error) {
	if len(jobId) == 0 {
		return nil, errors.New("unspecified argument 'jobId'")
	}
	log.Printf("Cancelling job %v", jobId)
	return commander.CombinedOutput("flink", "cancel", jobId)
}

func ListJobs() ([]byte, error) {
	return commander.CombinedOutput("flink", "list")
}

func Savepoint(jobId string, savepointTargetDir string) ([]byte, error) {
	if len(jobId) == 0 {
		return nil, errors.New("unspecified argument 'jobId'")
	}
	if len(savepointTargetDir) == 0 {
		return nil, errors.New("unspecified argument 'savepointTargetDir'")
	}
	log.Printf("Creating savepoint for job %v at targetDir %v", jobId, savepointTargetDir)
	return commander.CombinedOutput("flink", "savepoint", jobId, savepointTargetDir)
}
