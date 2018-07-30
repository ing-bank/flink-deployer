package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type CreateSavepointResponse struct {
	RequestId string `json:"request-id"`
}

type Status struct {
	Id string `json:"id"`
}

type Operation struct {
	Location string `json:"location"`
}

type CheckSavepointStatusResponse struct {
	Status    Status          `json:"status"`
	Operation json.RawMessage `json:"operation"`
}

func main() {
	r := strings.NewReader("{\"target-directory\": \"/data/flink/savepoints/testjob\",\"cancel-job\": false}")
	resp, err := http.Post("http://localhost:8081/jobs/c804e569341f3ea2321bb8e8bfee9873/savepoints", "application/json", r)

	if err != nil {
		fmt.Println("ERROR", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var createSavepointResponse CreateSavepointResponse
	json.Unmarshal(body, &createSavepointResponse)
	fmt.Println("Request ID:", createSavepointResponse.RequestId)

	time.Sleep(10 * time.Second)

	resp, err = http.Get(fmt.Sprintf("http://localhost:8081/jobs/c804e569341f3ea2321bb8e8bfee9873/savepoints/%s", createSavepointResponse.RequestId))
	if err != nil {
		fmt.Println("ERROR", err)
	}

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	var checkSavepointStatusResponse CheckSavepointStatusResponse
	json.Unmarshal(body, &checkSavepointStatusResponse)
	fmt.Println("RESPONSE: ", checkSavepointStatusResponse.Operation)

	fmt.Println(resp)

}
