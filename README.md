https://api.travis-ci.org/ing-bank/flink-deployer.svg?branch=master

# Flink-deployer

A Go command-line utility to facilitate deployments to Apache Flink.

Currently, it supports several features:

1. Listing jobs
2. Deploying a new job
3. Updating an existing job
4. Querying Flink queryable state

For a full overview of the commands and flags, run `flink-job-deployer help`

## Managing dependencies

This project uses [dep](https://github.com/golang/dep) to manage all project dependencies residing in the `vendor` folder. 

Run `dep status` to review the status of the included and most recent available depencencies.

## Build

Build from source for your current machine:

`go build ./cmd/cli`

Build from source for a specific machine architecture:

`env GOOS=linux GOARCH=amd64 go build ./cmd/cli`

Build the Docker container:

`go build ./cmd/cli`
`docker build -t com.ing/flink-job-deployer:latest .`

## Test

`go test`

### Test with coverage

`go test -coverprofile=cover.out && go tool cover`

## Copyright

All copyright of project flink-job-deployer are held by Marc Rooding and Niels Denissen, 2017-2018.
