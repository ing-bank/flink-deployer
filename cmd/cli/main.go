package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/afero"
	"github.com/urfave/cli"
)

var filesystem afero.Fs

func ListAction(c *cli.Context) error {
	out, err := ListJobs()

	log.Println(string(out))

	if err != nil {
		return err
	}

	return nil
}

func DeployAction(c *cli.Context) error {
	deploy := Deploy{}

	runArgs := c.String("run-args")
	if len(runArgs) > 0 {
		deploy.runArgs = runArgs
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
		deploy.localFilename = filename
	} else {
		deploy.remoteFilename = remoteFilename

		apiToken := c.String("api-token")
		if len(apiToken) > 0 {
			deploy.apiToken = apiToken
		}
	}

	jarArgs := c.String("jar-args")
	if len(jarArgs) > 0 {
		deploy.jarArgs = jarArgs
	}
	savepointPath := c.String("savepoint-path")
	if len(savepointPath) > 0 {
		deploy.savepointPath = savepointPath
	}
	deploy.allowNonRestorableState = c.Bool("allow-non-restorable-state")

	out, err := deploy.execute()

	log.Println(string(out))

	if err != nil {
		return cli.NewExitError(fmt.Sprintf("an error occurred: %v", err), -1)
	}

	return nil
}

func UpdateAction(c *cli.Context) error {
	update := UpdateJob{}

	jobNameBase := c.String("job-name-base")
	if len(jobNameBase) == 0 {
		return cli.NewExitError("unspecified flag 'job-name-base'", -1)
	} else {
		update.jobNameBase = jobNameBase
	}
	runArgs := c.String("run-args")
	if len(runArgs) > 0 {
		update.runArgs = runArgs
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
		update.localFilename = filename
	} else {
		update.remoteFilename = remoteFilename

		apiToken := c.String("api-token")
		if len(apiToken) > 0 {
			update.apiToken = apiToken
		}
	}
	jarArgs := c.String("jar-args")
	if len(jarArgs) > 0 {
		update.jarArgs = jarArgs
	}
	savepointDirectory := c.String("savepoint-dir")
	if len(savepointDirectory) > 0 {
		update.savepointDirectory = savepointDirectory
	}
	update.allowNonRestorableState = c.Bool("allow-non-restorable-state")

	out, err := update.execute()

	log.Println(string(out))

	if err != nil {
		return cli.NewExitError(fmt.Sprintf("an error occurred: %v", err), -1)
	}

	return nil
}

func QueryAction(c *cli.Context) error {
	query := Query{}

	jobName := c.String("job-name")
	if len(jobName) == 0 {
		return cli.NewExitError("unspecified flag 'job-name'", -1)
	} else {
		query.jobName = jobName
	}
	filename := c.String("file-name")
	if len(filename) == 0 {
		return cli.NewExitError("unspecified flag 'file-name'", -1)
	} else {
		query.filename = filename
	}
	mainClass := c.String("main-class")
	if len(mainClass) == 0 {
		return cli.NewExitError("unspecified flag 'main-class'", -1)
	} else {
		query.mainClass = mainClass
	}
	highAvailability := c.String("high-availability")
	if len(highAvailability) == 0 {
		return cli.NewExitError("unspecified flag 'high-availability'", -1)
	} else if highAvailability != "zookeeper" && highAvailability != "none" {
		return cli.NewExitError("unknown value for flag 'high-availability'. Allowed values: zookeeper/none", -1)
	} else {
		query.highAvailability = filename
	}
	zookeeperQuorum := c.String("zookeeper-quorum")
	if highAvailability == "zookeeper" && len(zookeeperQuorum) == 0 {
		return cli.NewExitError("unspecified flag 'zookeeper-quorum'", -1)
	} else {
		query.zookeeperQuorum = zookeeperQuorum
	}
	jobmanagerAddress := c.String("jobmanager-address")
	if highAvailability == "none" && len(jobmanagerAddress) == 0 {
		return cli.NewExitError("unspecified flag 'jobmanager-address'", -1)
	} else {
		query.jobManagerRPCAddress = jobmanagerAddress
	}
	jobmanagerPort := c.Int("jobmanager-port")
	if highAvailability == "none" && jobmanagerPort <= 0 {
		return cli.NewExitError("unspecified flag 'jobmanager-port'", -1)
	} else {
		query.jobManagerRPCPort = jobmanagerPort
	}

	out, err := query.execute()

	log.Println(string(out))

	if err != nil {
		return cli.NewExitError(fmt.Sprintf("an error occurred: %v", err), -1)
	}

	return nil
}

func main() {
	commander = RealCommander{}
	filesystem = afero.NewOsFs()

	app := cli.NewApp()
	app.Name = "flink-deployer"
	app.Description = "A Go command-line utility to facilitate deployments to Apache Flink"
	app.Version = "0.1.0"

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
					Usage: "The location of a remote job JAR file to be downloaded",
				},
				cli.StringFlag{
					Name:  "api-token, at",
					Usage: "The API token for the remote address of the a remote file",
				},
				cli.StringFlag{
					Name:  "run-args, ra",
					Usage: "The arguments to pass to the 'flink run' command",
				},
				cli.StringFlag{
					Name:  "jar-args, ja",
					Usage: "The arguments to pass to the jar execution",
				},
				cli.StringFlag{
					Name:  "savepoint-path, sp",
					Usage: "The path to the savepoint to restore from",
				},
				cli.BoolFlag{
					Name:  "allow-non-restorable-state, anrs",
					Usage: "Allow non restored savepoint state in case an operator has been removed from the job.",
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
					Usage: "The location of a remote job JAR file to be downloaded",
				},
				cli.StringFlag{
					Name:  "api-token, at",
					Usage: "The API token for the remote address of the a remote file",
				},
				cli.StringFlag{
					Name:  "run-args, ra",
					Usage: "The arguments to pass to the 'flink run' command",
				},
				cli.StringFlag{
					Name:  "jar-args, ja",
					Usage: "The arguments to pass to the jar execution",
				},
				cli.StringFlag{
					Name:  "savepoint-dir, sd",
					Usage: "The path to the directory where Flink stores all savepoints",
				},
				cli.BoolFlag{
					Name:  "allow-non-restorable-state, anrs",
					Usage: "The savepoint directory to restore a savepoint from",
				},
			},
			Action: UpdateAction,
		},
		{
			Name:    "query",
			Aliases: []string{"q"},
			Usage:   "run a query against a job's state",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "job-name, jn",
					Usage: "The name of the job to update",
				},
				cli.StringFlag{
					Name:  "file-name, fn",
					Usage: "The complete name of the job JAR file",
				},
				cli.StringFlag{
					Name:  "remote-file-name, rfn",
					Usage: "The location of a remote job JAR file to be downloaded",
				},
				cli.StringFlag{
					Name:  "api-token, at",
					Usage: "The API token for the remote address of the a remote file",
				},
				cli.StringFlag{
					Name:  "main-class, mc",
					Usage: "The package and class name of the main class",
				},
				cli.StringFlag{
					Name:  "high-availability, ha",
					Usage: "Which high availability mode to use (zookeeper/none)",
				},
				cli.StringFlag{
					Name:  "zookeeper-quorum, zq",
					Usage: "Comma separated list of zookeeper nodes (host:port,host:port)",
				},
				cli.StringFlag{
					Name:  "jobmanager-address, ja",
					Usage: "The Job Manager RPC address to use when high availability is none",
				},
				cli.IntFlag{
					Name:  "jobmanager-port, jp",
					Usage: "The Job Manager RPC port to use when high availability is none",
				},
			},
			Action: QueryAction,
		},
	}

	app.Run(os.Args)
}
