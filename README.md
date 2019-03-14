[![Build Status](https://travis-ci.org/ing-bank/flink-deployer.svg?branch=master)](https://travis-ci.org/ing-bank/flink-deployer)
[![codecov.io](http://codecov.io/github/ing-bank/flink-deployer/coverage.svg?branch=master)](https://codecov.io/gh/ing-bank/flink-deployer?branch=master)
[![](https://images.microbadger.com/badges/image/nielsdenissen/flink-deployer:master.svg)](https://hub.docker.com/r/nielsdenissen/flink-deployer/)
[![](https://images.microbadger.com/badges/version/nielsdenissen/flink-deployer:master.svg)](https://hub.docker.com/r/nielsdenissen/flink-deployer/)
[![Gitter chat](https://badges.gitter.im/gitterHQ/gitter.png)](https://gitter.im/Flink-Deployer/Lobby)

# Flink-deployer

A Go command-line utility to facilitate deployments to Apache Flink.

Currently, it supports several features:

1. Listing jobs
2. Deploying a new job
3. Updating an existing job
4. Querying Flink queryable state

For a full overview of the commands and flags, run `flink-job-deployer help`

## <a name="howtorunlocally"></a>How to run locally

To be able to test the deployer locally, follow these steps:

1. Build the CLI tool docker image: `docker-compose build deployer`
2. ***optional***: `cd flink-sample-job; sbt clean assembly; cd ..` (Builds a jar with small stateful test job)
3. `docker-compose up -d jobmanager taskmanager` (start a Flink job- and taskmanager)
4. `docker-compose run deployer help` (run the Flink deployer with argument `help`)

Repeat step 3 with any commands you'd like to try. 

### Run a sample job
Provided you ran step 1 of the above guide, a jar with a sample Flink job is available in the deployer. It will be mounted in the deployer container at the following path:

    /tmp/flink-sample-job/flink-stateful-wordcount-assembly-0.jar

To deploy it you can simply run (it's the default command specified in the `docker-compose.yml`): 

```bash
docker-compose run deployer
```

This will print a simple word count to the output console, you can view it by checking the logs of the taskmanager as follows:

```bash
docker-compose logs -f taskmanager
```

If all went well you should see the word counter continue with where it was.

A list of some example commands to run can be found [here](./docs/example-commands.md).

## Authentication

Apache Flink doesn't support any Web UI authentication out of the box. One of the custom approaches is using NGINX in front of Flink to protect the user interface. With NGINX, there are again a lot of different ways to add that authentication layer. To support the most basic one, we've added support for using Basic Authentication.

You can inject the `FLINK_BASIC_AUTH_USERNAME` and `FLINK_BASIC_AUTH_PASSWORD` environment variables to configure basic authentication.

## Supported environment variables

* FLINK_BASE_URL: Base Url to Flink's API (**required**, e.g. http://jobmanageraddress:8081/)
* FLINK_BASIC_AUTH_USERNAME: Basic authentication username used for authenticating to Flink 
* FLINK_BASIC_AUTH_PASSWORD: Basic authentication password used for authenticating to Flink 
* FLINK_API_TIMEOUT_SECONDS: Number of seconds until requests to the Flink API time out (e.g. 10)

# Development

## Managing dependencies

This project uses [dep](https://github.com/golang/dep) to manage all project dependencies residing in the `vendor` folder. 

Run `dep status` to review the status of the included and most recent available depencencies.

## Build

Build from source for your current machine:

```bash
go build ./cmd/cli
```

Build from source for a specific machine architecture:

```bash
env GOOS=linux GOARCH=amd64 go build ./cmd/cli
```

Build the Docker container locally to test CLI tool:

```bash
docker-compose build deployer
```

## Test

```bash
go test ./cmd/cli ./cmd/cli/flink ./cmd/cli/operations
```

Or with coverage:

```bash
sh test-with-coverage.sh
```

# Docker

A docker image for this repo is available from the docker hub: `nielsdenissen/flink-deployer`

The image expects the following env vars:

```bash
FLINK_BASE_URL=http://localhost:8080
```


# Kubernetes

When running in Kubernetes (or Openshift), you'll have to deploy the container to the cluster. A reason for this is Flink will try to reroute you to the internal Kubernetes address of the cluster, which doesn't resolve from outside. Besides that it'll give you the necessary access to the stored savepoints when you're using persistent volumes to store those.

This section is aimed at providing you with a quick getting started guide to deploy our container to Kubernetes. There are a few steps we'll need to take which we describe below:

## 0. Run a kubernetes cluster

If you don't have a kubernetes cluster readily available, you can quickly get started by setting up a [minikube cluster](https://kubernetes.io/docs/setup/minikube/).

    minikube start

## 1. Setup a Flink cluster in Kubernetes

Flink has a guide on how to run a cluster in Kubernetes, you can find it [here](https://ci.apache.org/projects/flink/flink-docs-stable/ops/deployment/kubernetes.html).

>If you're using Minikube, be sure to pull the images that flink uses in their deploy configurations locally first. Otherwise minikube will not be able to find them. So perform a `docker pull flink:latest` on your host.

## 2. Add the test jar (or your own job you want to run) to the deployer image

We now need to package the jar into the container so we can deploy it in Kubernetes. There are other ways around this like storing the jar on a Persistent Volume or downloading it at runtime inside the container. This is the easiest getting started though and still the technique we use.

To build the container with the jar packaged you can use the `Dockerfile-including-sample-job`. Be sure to have create the jar for the test job in case you want to use it. See step 2 in the [How to run locally](#howtorunlocally) section.
Run this from the root of this repository:
 
    docker build -t flinkdeployerstatefulwordcount:test -f Dockerfile-including-sample-job .

## 3. Run the deployer in Kubernetes

In this example we're going to show how you can do a simple deploy of the sample-job in this project to the cluster. For this we need a yaml that specifies what to do to Kubernetes. Here's an example of how such a kubernetes yaml could look like:

```yaml
apiVersion: v1
kind: Pod
metadata:
    generateName: "flink-stateful-wordcount-deployer-"
spec:
    dnsPolicy: ClusterFirst
    restartPolicy: OnFailure
    containers:
    -   name: "flink-stateful-wordcount-deployer"
        image: "flinkdeployerstatefulwordcount:test"
        args:
        - "deploy"
        - "--file-name"
        - "/tmp/flink-stateful-wordcount-assembly-0.jar"
        - "--entry-class"
        - "WordCountStateful"
        - "--parallelism"
        - "2"
        - "--program-args"
        - "--intervalMs 1000"
        imagePullPolicy: Never
        env:
        -   name: FLINK_BASE_URL
            value: "http://flink-jobmanager:8081"
```

Go to Kubernetes, click the `Create +` button and copy paste the above YAML. This should trigger a POD to be deployed that runs once and stops after deploying the sample job to the Flink cluster running in Kubernetes.

**MINIKUBE USERS**: In order to use local images with Minikube (so images on your local docker installation instead of dockerHub), you need to perform the following steps:
- Point minikube to your local docker: `eval $(minikube docker-env)` (See [this guide](https://blogmilind.wordpress.com/2018/01/30/running-local-docker-images-in-kubernetes/) for more info)
- Rebuild the image as done in step 2 of this guide.
- The `imagePullPolicy` in the yaml above must be set to `Never`.

## 4. Attach Persistent Volumes to all Flink containers

This step we won't outline completely, as it's a bit more involved for a getting started guide. In order to recover jobs from savepoints, you'll need to have a Persistent Volume shared among all Flink nodes and the deployer. You'll need this in any case if you want to persistent and thus not lose any data in your Flink cluster running in Kubernetes.
After creating a Persistent Volume and hooking it up to the existing Flink containers, you'll need to add something like the following to the YAML of the deployer (besides of course change the command to for instance `update`):

```yaml
        volumeMounts:
        -   name: flink-data
            mountPath: "/data/flink"
    volumes:
    -   name: flink-data
        persistentVolumeClaim:
            claimName: "PVC_FLINK"
```

The directory you put in your Persistent Volume should be the directory to which Flink stores it's savepoints.

# Copyright

All copyright of project flink-job-deployer are held by Marc Rooding and Niels Denissen, 2017-2018.
