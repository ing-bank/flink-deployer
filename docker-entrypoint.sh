#!/usr/bin/env bash

# Function that modifies the flink configuration file according to the job manager deployment strategy
function generateConfigurationFile {
    envsubst < /conf/flink-conf-template.yaml > ${FLINK_HOME}/conf/flink-conf.yaml

    if [ ${HIGH_AVAILABILITY} == "zookeeper" ]; then
        sed -i -e 's/jobmanager.rpc.address/#jobmanager.rpc.address/g' ${FLINK_HOME}/conf/flink-conf.yaml
        sed -i -e 's/jobmanager.rpc.port/#jobmanager.rpc.port/g' ${FLINK_HOME}/conf/flink-conf.yaml
    fi

    echo "config file: " && grep '^[^\n#]' "$FLINK_HOME/conf/flink-conf.yaml"
}

generateConfigurationFile

/flink-deployer/cli $FLINK_DEPLOY_OPERATION -j "$FLINK_JOB_NAME" --rf "$FLINK_REMOTE_FILE" --at "$GITLAB_API_TOKEN" --ra "$FLINK_RUN_ARGS" --ja "$FLINK_JAR_ARGS" --sd "$FLINK_SAVEPOINT_DIR"
