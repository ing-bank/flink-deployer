#!/usr/bin/env bash

# Function that modifies the flink configuration file according to the job manager deployment strategy
function generateConfigurationFile {
    envsubst < /flink-deployer/conf/flink-conf-template.yaml > ${FLINK_HOME}/conf/flink-conf.yaml

    if [ ${HIGH_AVAILABILITY} == "zookeeper" ]; then
        sed -i -e 's/jobmanager.rpc.address/#jobmanager.rpc.address/g' ${FLINK_HOME}/conf/flink-conf.yaml
        sed -i -e 's/jobmanager.rpc.port/#jobmanager.rpc.port/g' ${FLINK_HOME}/conf/flink-conf.yaml
    fi

    echo "config file: " && grep '^[^\n#]' "$FLINK_HOME/conf/flink-conf.yaml"
}

generateConfigurationFile

echo "Running cli with command: $@"
exec /flink-deployer/cli "$@"
