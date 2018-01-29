package main

import (
	"os/exec"
	"os"
	"strconv"
	"fmt"
	"testing"
)

/*
 * Commander mocking
 */
var mockedExitStatus = 0
var mockedStdout string
var receivedArguments []string

type TestCommander struct{}

func (c TestCommander) CombinedOutput(command string, args ...string) ([]byte, error) {
	receivedArguments = args
	cs := []string{"-test.run=TestExecCommandHelper", "--"}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	es := strconv.Itoa(mockedExitStatus)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1",
		"GO_TEST_STDOUT=" + mockedStdout,
		"GO_TEXT_EXIT_STATUS=" + es}
	out, err := cmd.CombinedOutput()
	return out, err
}

func TestExecCommandHelper(*testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	fmt.Fprintf(os.Stdout, os.Getenv("GO_TEST_STDOUT"))
	i, _ := strconv.Atoi(os.Getenv("GO_TEXT_EXIT_STATUS"))
	os.Exit(i)
}
