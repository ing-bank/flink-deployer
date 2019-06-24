package flink

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

type runJarRequest struct {
	EntryClass            string `json:"entryClass"`
	ProgramArgs           string `json:"programArgs"`
	Parallelism           int    `json:"parallelism"`
	AllowNonRestoredState bool   `json:"allowNonRestoredState"`
	SavepointPath         string `json:"savepointPath"`
}

// RunJar executes a specific JAR file with the supplied parameters on the Flink cluster
func (c FlinkRestClient) RunJar(jarID string, entryClass string, jarArgs []string, parallelism int, savepointPath string, allowNonRestoredState bool) error {
	runJarRequest := runJarRequest{
		EntryClass:            entryClass,
		ProgramArgs:           strings.Join(jarArgs, " "),
		Parallelism:           parallelism,
		AllowNonRestoredState: allowNonRestoredState,
		SavepointPath:         savepointPath,
	}
	reqBody := new(bytes.Buffer)
	json.NewEncoder(reqBody).Encode(runJarRequest)

	req, err := c.newRequest("POST", c.constructURL(fmt.Sprintf("jars/%v/run", jarID)), reqBody)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := c.Client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("Unexpected response status %v with body %v", res.StatusCode, string(resBody[:]))
	}

	return nil
}
