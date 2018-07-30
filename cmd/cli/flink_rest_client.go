package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

var netClient = &http.Client{
	Timeout: time.Second * 2,
}

type FlinkRestClient struct {
	host string
	port int
}

type Job struct {
	id     string `json:"id"`
	name   string `json:"name"`
	status string `json:"state"`
}

type retrieveJobsResponse struct {
	jobs []Job `json:"jobs"`
}

func (c FlinkRestClient) retrieveJobs() ([]Job, error) {
	res, err := netClient.Get(c.constructURL("jobs/overview"))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []Job{}, err
	}

	retrieveJobsResponse := retrieveJobsResponse{}
	err = json.Unmarshal(body, &retrieveJobsResponse)
	if err != nil {
		return []Job{}, errors.New("Unable to parse API response as valid JSON")
	}

	return retrieveJobsResponse.jobs, nil
}

type Jar struct {
	id   string `json:"id"`
	name string `json:"name"`
}

type retrieveJarsResponse struct {
	address string `json:"address"`
	files   []Jar  `json:"files"`
}

func (c FlinkRestClient) retrieveJars() ([]Jar, error) {
	res, err := netClient.Get(c.constructURL("jars"))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []Jar{}, err
	}

	retrieveJarsResponse := retrieveJarsResponse{}
	err = json.Unmarshal(body, &retrieveJarsResponse)
	if err != nil {
		return []Jar{}, errors.New("Unable to parse API response as valid JSON")
	}

	return retrieveJarsResponse.files, nil
}

type uploadJarResponse struct {
	filename string `json:"filename"`
	status   string `json:"status"`
}

func (c FlinkRestClient) uploadJar(filename string) (uploadJarResponse, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	fileWriter, err := bodyWriter.CreateFormFile("jarfile", filename)
	if err != nil {
		return uploadJarResponse{}, err
	}

	fh, err := os.Open(filename)
	if err != nil {
		return uploadJarResponse{}, err
	}
	defer fh.Close()

	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return uploadJarResponse{}, err
	}

	bodyWriter.Close()

	res, err := netClient.Post(c.constructURL("jars/upload"), "application/x-java-archive", bodyBuf)
	if err != nil {
		return uploadJarResponse{}, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return uploadJarResponse{}, err
	}

	response := uploadJarResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return uploadJarResponse{}, errors.New("Unable to parse API response as valid JSON")
	}

	return response, nil
}

func (c FlinkRestClient) runJar(jarID string, jarArgs string, savepointPath string, allowNonRestorableState bool) error {
	res, err := netClient.Get(c.constructURL(fmt.Sprintf("jars/%v/run?program-args=%v&allowNonRestoredState=%v&savepointPath=%v", jarID, jarArgs, allowNonRestorableState, savepointPath)))
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("Unexpected response status: %v", res.Status)
	}

	return nil
}

func (c FlinkRestClient) constructURL(path string) string {
	return fmt.Sprintf("http://%v:%d/%v", c.host, c.port, path)
}
