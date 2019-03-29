package flink

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// UploadJarResponse represents the response body
// used by the upload JAR API
type UploadJarResponse struct {
	Filename string `json:"filename"`
	Status   string `json:"status"`
}

func (c FlinkRestClient) constructUploadJarRequest(filename string, url string) (*http.Response, error) {
	buffer := &bytes.Buffer{}
	writer := multipart.NewWriter(buffer)

	file, err := os.Open(filename)
	if err != nil {
		return &http.Response{}, err
	}
	defer file.Close()

	part, err := writer.CreateFormFile("jarfile", filepath.Base(filename))
	if err != nil {
		return &http.Response{}, err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return &http.Response{}, err
	}

	contentType := writer.FormDataContentType()
	writer.Close()

	req, err := c.newRequest("POST", url, buffer)
	if err != nil {
		return &http.Response{}, err
	}
	req.Header.Set("Content-Type", contentType)

	return c.Client.Do(req)
}

// UploadJar allows for uploading a JAR file to the Flink cluster
func (c FlinkRestClient) UploadJar(filename string) (UploadJarResponse, error) {
	res, err := c.constructUploadJarRequest(filename, c.constructURL("jars/upload"))
	if err != nil {
		return UploadJarResponse{}, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return UploadJarResponse{}, err
	}

	if res.StatusCode != 200 {
		return UploadJarResponse{}, fmt.Errorf("Unexpected response status %v with body %v", res.StatusCode, string(body[:]))
	}

	response := UploadJarResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return UploadJarResponse{}, fmt.Errorf("Unable to parse API response as valid JSON: %v", string(body[:]))
	}

	return response, nil
}
