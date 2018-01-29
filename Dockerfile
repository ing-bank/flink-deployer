FROM flink:1.3.2

# Required to run our cli binary
RUN apt-get update && \
    apt-get install -y musl gettext procps vim

ENV HIGH_AVAILABILITY=none \
    JOB_MANAGER_WEB_ADDRESS=0.0.0.0 \
    JOB_MANAGER_WEB_PORT=8081 \
    JOB_MANAGER_RPC_ADDRESS=localhost \
    JOB_MANAGER_RPC_PORT=6123 \
    BLOB_SERVER_PORT=6124 \
    QUERY_SERVER_PORT=6125

COPY docker-entrypoint.sh /flink-deployer/docker-entrypoint.sh
COPY conf /flink-deployer/
COPY cli /flink-deployer/cli

RUN chmod -R 775 /flink-deployer ${FLINK_HOME}

WORKDIR /flink-deployer

VOLUME [ "/data/flink" ]
ENTRYPOINT [ "/flink-deployer/docker-entrypoint.sh" ]
