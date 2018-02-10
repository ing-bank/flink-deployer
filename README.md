[![Build Status](https://travis-ci.org/ing-bank/flink-deployer.svg?branch=master)](https://travis-ci.org/ing-bank/flink-deployer)
[![codecov.io](http://codecov.io/github/ing-bank/flink-deployer/coverage.svg?branch=master)](https://codecov.io/gh/ing-bank/flink-deployer?branch=master)
[![](https://images.microbadger.com/badges/image/nielsdenissen/flink-deployer:master.svg)](https://microbadger.com/images/nielsdenissen/flink-deployer:master)
[![](https://images.microbadger.com/badges/version/nielsdenissen/flink-deployer:master.svg)](https://microbadger.com/images/nielsdenissen/flink-deployer:master)

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

# Docker image

A docker image for this repo is available from the docker hub: `nielsdenissen/flink-deployer`

# Copyright

All copyright of project flink-job-deployer are held by Marc Rooding and Niels Denissen, 2017-2018.
