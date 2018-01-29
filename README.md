# Flink-job-deployer

A Go command-line utility to facilitate deployments to Apache Flink.

Currently, it supports several features:

1. Listing jobs
2. Deploying a new job
3. Updating an existing job
4. Querying Flink queryable state

For a full overview of the commands and flags, run `flink-job-deployer help`

## Managing dependencies

This project uses [dep](https://github.com/golang/dep) to manage all project dependencies residing in the `vendor` folder. 

## Build

`env GOOS=linux GOARCH=amd64 go build`

## Test

`go test`

### Test with coverage

`go test -coverprofile=cover.out && go tool cover`
