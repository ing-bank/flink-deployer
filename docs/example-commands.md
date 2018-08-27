# Example Commands

This file shows some commands that you can perform with the Flink Deployer.

1. Deploy a new job

```bash
docker-compose run deployer deploy \
    --file-name "/tmp/flink-stateful-wordcount-assembly-0.jar" \
    --entry-class "WordCountStateful" \
    --parallelism "2" \
    --program-args "--intervalMs 1000"
```

2. List running jobs

```bash
docker-compose run deployer list
```

3. Upgrade a running job

```bash
docker-compose run deployer update \
    --job-name-base "Windowed WordCount" \
    --file-name "/tmp/flink-stateful-wordcount-assembly-0.jar" \
    --entry-class "WordCountStateful" \
    --parallelism "2" \
    --program-args "--intervalMs 1000" \
    --savepoint-dir "/data/flink"
```

4. Start a job from a specific savepoint
Ensure you've run steps 1 and 3. This will have created a savepoint. Find the location of that savepoint and put this in the placeholder below:

```bash
docker-compose run deployer deploy \
    --file-name "/tmp/flink-stateful-wordcount-assembly-0.jar" \
    --entry-class "WordCountStateful" \
    --parallelism "2" \
    --program-args "--intervalMs 1000" \
    --savepoint-path "/data/flink/[SAVEPOINT_LOC_HERE]"
```

5. Start a job from the latest savepoint in a specified savepoint directory

Ensure you've run steps 1 and 3. This will have created a savepoint.

```bash
docker-compose run deployer deploy \
    --file-name "/tmp/flink-stateful-wordcount-assembly-0.jar" \
    --entry-class "WordCountStateful" \
    --parallelism "2" \
    --program-args "--intervalMs 1000" \
    --savepoint-dir "/data/flink"
```