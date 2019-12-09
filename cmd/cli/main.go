package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/ing-bank/flink-deployer/cmd/cli/flink"
	"github.com/ing-bank/flink-deployer/cmd/cli/operations"
	"github.com/spf13/afero"
	"github.com/urfave/cli"
)

var filesystem afero.Fs
var operator operations.Operator

// ListAction executes the CLI list command
func ListAction(c *cli.Context) error {
	jobs, err := operator.RetrieveJobs()
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("failed to list jobs: %v", err), -1)
	}

	if len(jobs) == 0 {
		log.Println("No running jobs found")
	} else {
		for _, job := range jobs {
			log.Printf("Job %v (%v) with status: %v", job.Name, job.ID, job.Status)
		}
	}

	return nil
}

// DeployAction executes the CLI deploy command
func DeployAction(c *cli.Context) error {
	deploy := operations.Deploy{}

	filename := c.String("file-name")
	remoteFilename := c.String("remote-file-name")
	if len(filename) == 0 && len(remoteFilename) == 0 {
		return cli.NewExitError("both flags 'file-name' and 'remote-file-name' unspecified", -1)
	}
	if len(filename) > 0 && len(remoteFilename) > 0 {
		return cli.NewExitError("both flags 'file-name' and 'remote-file-name' specified, only one allowed", -1)
	}

	if len(filename) > 0 {
		deploy.LocalFilename = filename
	} else {
		deploy.RemoteFilename = remoteFilename

		apiToken := c.String("api-token")
		if len(apiToken) > 0 {
			deploy.APIToken = apiToken
		}
	}

	entryClass := c.String("entry-class")
	if len(entryClass) > 0 {
		deploy.EntryClass = entryClass
	}

	parallelism := c.Int("parallelism")
	if parallelism != 0 {
		deploy.Parallelism = parallelism
	} else {
		deploy.Parallelism = 1
	}

	programArgs := c.StringSlice("program-args")
	if len(programArgs) > 0 {
		deploy.ProgramArgs = programArgs
	}

	savepointDir := c.String("savepoint-dir")
	savepointPath := c.String("savepoint-path")
	if len(savepointDir) > 0 && len(savepointPath) > 0 {
		return cli.NewExitError("both flags 'savepoint-dir' and 'savepoint-path' specified, only one allowed", -1)
	}
	if len(savepointDir) > 0 {
		deploy.SavepointDir = savepointDir
	}
	if len(savepointPath) > 0 {
		deploy.SavepointPath = savepointPath
	}

	deploy.AllowNonRestoredState = c.Bool("allow-non-restored-state")

	err := operator.Deploy(deploy)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("an error occurred: %v", err), -1)
	}

	log.Println("Job started successfully")

	return nil
}

// UpdateAction executes the CLI update command
func UpdateAction(c *cli.Context) error {
	update := operations.UpdateJob{}

	jobNameBase := c.String("job-name-base")
	if len(jobNameBase) != 0 {
		update.JobNameBase = jobNameBase
	} else {
		return cli.NewExitError("unspecified flag 'job-name-base'", -1)
	}

	filename := c.String("file-name")
	remoteFilename := c.String("remote-file-name")
	if len(filename) == 0 && len(remoteFilename) == 0 {
		return cli.NewExitError("both flags 'file-name' and 'remote-file-name' unspecified", -1)
	}
	if len(filename) > 0 && len(remoteFilename) > 0 {
		return cli.NewExitError("both flags 'file-name' and 'remote-file-name' specified, only one allowed", -1)
	}
	if len(filename) > 0 {
		update.LocalFilename = filename
	} else {
		update.RemoteFilename = remoteFilename

		apiToken := c.String("api-token")
		if len(apiToken) > 0 {
			update.APIToken = apiToken
		}
	}

	entryClass := c.String("entry-class")
	if len(entryClass) > 0 {
		update.EntryClass = entryClass
	}

	parallelism := c.Int("parallelism")
	if parallelism != 0 {
		update.Parallelism = parallelism
	} else {
		update.Parallelism = 1
	}

	programArgs := c.StringSlice("program-args")
	if len(programArgs) > 0 {
		update.ProgramArgs = programArgs
	}

	savepointDir := c.String("savepoint-dir")
	if len(savepointDir) != 0 {
		update.SavepointDir = savepointDir
	} else {
		return cli.NewExitError("unspecified flag 'savepoint-dir'", -1)
	}

	update.AllowNonRestoredState = c.Bool("allow-non-restored-state")

	update.FallbackToDeploy = c.Bool("fallback-to-deploy")

	err := operator.Update(update)

	if err != nil {
		return cli.NewExitError(fmt.Sprintf("an error occurred: %v", err), -1)
	}

	log.Println("Job successfully updated")

	return nil
}

// TerminateAction executes the CLI terminate command
func TerminateAction(c *cli.Context) error {
	terminate := operations.TerminateJob{}

	jobNameBase := c.String("job-name-base")
	if len(jobNameBase) == 0 {
		return cli.NewExitError("unspecified flag 'job-name-base'", -1)
	}
	terminate.JobNameBase = jobNameBase

	mode := c.String("mode")
	if len(mode) > 0 && mode != "cancel" && mode != "stop" {
		return cli.NewExitError("unknown value for 'mode', only 'cancel' and 'stop' are supported", -1)
	}
	terminate.Mode = mode

	err := operator.Terminate(terminate)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("an error occurred: %v", err), -1)
	}

	log.Println("Job successfully terminated")

	return nil
}

func getAPITimeoutSeconds() (int64, error) {
	if len(os.Getenv("FLINK_API_TIMEOUT_SECONDS")) > 0 {
		return strconv.ParseInt(os.Getenv("FLINK_API_TIMEOUT_SECONDS"), 10, 64)
	}
	return int64(10), nil
}

func main() {
	flinkBaseURL := os.Getenv("FLINK_BASE_URL")
	if len(flinkBaseURL) == 0 {
		log.Fatal("`FLINK_BASE_URL` environment variable not found")
		os.Exit(1)
	}

	flinkBasicAuthUsername := os.Getenv("FLINK_BASIC_AUTH_USERNAME")
	flinkBasicAuthPassword := os.Getenv("FLINK_BASIC_AUTH_PASSWORD")

	flinkAPITimeoutSeconds, err := getAPITimeoutSeconds()
	if err != nil {
		log.Fatalf("`FLINK_API_TIMEOUT_SECONDS=%v` environment variable could not be parsed to an integer", os.Getenv("FLINK_API_TIMEOUT_SECONDS"))
		os.Exit(1)
	}

	client := retryablehttp.NewClient()
	client.HTTPClient = &http.Client{
		Timeout: time.Second * time.Duration(flinkAPITimeoutSeconds),
	}

	operator = operations.RealOperator{
		Filesystem: afero.NewOsFs(),
		FlinkRestAPI: flink.FlinkRestClient{
			BaseURL:           flinkBaseURL,
			BasicAuthUsername: flinkBasicAuthUsername,
			BasicAuthPassword: flinkBasicAuthPassword,
			Client:            client,
		},
	}

	app := cli.NewApp()
	app.Name = "Flink Deployer"
	app.Description = "A Go command-line utility to facilitate deployments to Apache Flink"
	app.Version = "1.3.0"

	app.Commands = []cli.Command{
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "list the jobs running on the job manager",
			Action:  ListAction,
		},
		{
			Name:    "deploy",
			Aliases: []string{"d"},
			Usage:   "deploy the JAR to the job manager",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "file-name, fn",
					Usage: "The complete name of the job JAR file",
				},
				cli.StringFlag{
					Name:  "remote-file-name, rfn",
					Usage: "The location of a GitLab job JAR file to be downloaded",
				},
				cli.StringFlag{
					Name:  "api-token, at",
					Usage: "The GitLab API token for the remote address of the a remote file",
				},
				cli.StringFlag{
					Name:  "entry-class, ec",
					Usage: "The entry class name that contains the main method",
				},
				cli.StringFlag{
					Name:  "parallelism, p",
					Usage: "The parallelism count",
				},
				cli.StringSliceFlag{
					Name:  "program-args, pa",
					Usage: "The arguments to pass to the program execution. This flag may be repeated to provide multiple arguments",
				},
				cli.StringFlag{
					Name:  "savepoint-dir, sd",
					Usage: "The path to the directory that contains the savepoints",
				},
				cli.StringFlag{
					Name:  "savepoint-path, sp",
					Usage: "The path to the savepoint to restore from",
				},
				cli.BoolFlag{
					Name:  "allow-non-restored-state, anrs",
					Usage: "Allow the job to run if the state cannot be restored",
				},
			},
			Action: DeployAction,
		},
		{
			Name:    "update",
			Aliases: []string{"u"},
			Usage:   "Update a running job by creating a savepoint, stopping the job and deploying the new version",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "job-name-base, jnb",
					Usage: "The base name of the job to update",
				},
				cli.StringFlag{
					Name:  "file-name, fn",
					Usage: "The complete name of the job JAR file",
				},
				cli.StringFlag{
					Name:  "remote-file-name, rfn",
					Usage: "The location of a GitLab job JAR file to be downloaded",
				},
				cli.StringFlag{
					Name:  "api-token, at",
					Usage: "The GitLab API token for the remote address of the a remote file",
				},
				cli.StringFlag{
					Name:  "entry-class, ec",
					Usage: "The entry class name that contains the main method",
				},
				cli.StringFlag{
					Name:  "parallelism, p",
					Usage: "The parallelism count",
				},
				cli.StringSliceFlag{
					Name:  "program-args, pa",
					Usage: "The arguments to pass to the program execution",
				},
				cli.StringFlag{
					Name:  "savepoint-dir, sd",
					Usage: "The path to the directory that contains the savepoints",
				},
				cli.BoolFlag{
					Name:  "allow-non-restored-state, anrs",
					Usage: "Allow the job to run if the state cannot be restored",
				},
				cli.BoolFlag{
					Name:  "fallback-to-deploy, fbd",
					Usage: "Continue to deploy the job if no running instance of the job is found",
				},
			},
			Action: UpdateAction,
		},
		{
			Name:    "terminate",
			Aliases: []string{"t"},
			Usage:   "Terminate a running job",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "job-name-base, jnb",
					Usage: "The base name of the job to update",
				},
				cli.StringFlag{
					Name:  "mode, m",
					Usage: "The mode to terminate a running job, cancel and stop supported",
				},
			},
			Action: TerminateAction,
		},
	}

	app.Run(os.Args)
}
