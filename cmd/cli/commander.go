package main

import "os/exec"

var commander Commander

type Commander interface {
	CombinedOutput(string, ...string) ([]byte, error)
}

type RealCommander struct{}

func (c RealCommander) CombinedOutput(command string, args ...string) ([]byte, error) {
	return exec.Command(command, args...).CombinedOutput()
}
