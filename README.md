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

Or with coverage:

`go test -coverprofile=cover.out && go tool cover`

# Docker

A docker image for this repo is available from the docker hub: `nielsdenissen/flink-deployer`

# Kubernetes

When running in Kubernetes (or Openshift), you'll have to deploy the container to the cluster. A reason for this is Flink will try to reroute you to the internal Kubernetes address of the cluster, which doesn't resolve from outside. Besides that it'll give you the necessary access to the stored savepoints when you're using persistent volumes to store those.

Here's an example of how such a kubernetes yaml could look like:

```yaml
    apiVersion: v1
    kind: Template
    objects:
    -   apiVersion: v1
        kind: Pod
        metadata:
            generateName: "flink-${FLINK_JOB_ID}-deployer-"
        spec:
            dnsPolicy: ClusterFirst
            restartPolicy: OnFailure
            containers:
            -   name: "flink-${FLINK_JOB_ID}-deployer"
                image: "nielsdenissen/flink-deployer"
                args:
                - "update"
                - "-job-name-base"
                - "$(FLINK_JOB_NAME_BASE)"
                - "-file-name"
                - "/tmp/YOUR-FLINK-JAR.jar"
                - "-run-args"
                - "-p 2 -d -c $(MAIN_CLASS_NAME)"
                - "-jar-args"
                - "--kafka.bootstrapServers $(KAFKA_BOOTSTRAPSERVERS)"
                - "-savepoint-dir"
                - "/data/flink/savepoints/$(FLINK_JOB_ID)"
                imagePullPolicy: Always
                env:
                -   name: FLINK_JOB_NAME_BASE
                    value: "${FLINK_JOB_NAME_BASE}"
                -   name: JOB_MANAGER_RPC_ADDRESS
                    value: "jobmanager"
                -   name: JOB_MANAGER_RPC_PORT
                    value: "8081"
                -   name: HIGH_AVAILABILITY
                    value: "zookeeper"
                -   name: ZOOKEEPER_QUORUM
                    value: "zookeeper:2181"
                -   name: KAFKA_BOOTSTRAPSERVERS
                    value: "kafka:9092"
                -   name: MAIN_CLASS_NAME
                    value: "${MAIN_CLASS_NAME}"
                -   name: FLINK_JOB_ID
                    value: "${FLINK_JOB_ID}"
                volumeMounts:
                -   name: flink-data
                    mountPath: "/data/flink"
            volumes:
            -   name: flink-data
                persistentVolumeClaim:
                    claimName: "${PVC_FLINK}"
    parameters:
    -   name: FLINK_JOB_ID
        description: The ID to use for pod name and savepoint directory
    -   name: FLINK_JOB_NAME_BASE
        description: The job name base (you can append a version number behind this base in your actual job name)
    -   name: MAIN_CLASS_NAME
        description: Name of the main class to be run in the JAR
    -   name: PVC_FLINK
        description: The persistent volume claim name for flink.
```

# Copyright

All copyright of project flink-job-deployer are held by Marc Rooding and Niels Denissen, 2017-2018.
