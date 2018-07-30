FROM alpine:3.8

ENV HIGH_AVAILABILITY=none \
    JOB_MANAGER_RPC_ADDRESS=localhost \
    JOB_MANAGER_RPC_PORT=6123 \
    JOB_MANAGER_WEB_ADDRESS=0.0.0.0 \
    JOB_MANAGER_WEB_PORT=8081

COPY cli /flink-deployer/cli

WORKDIR /flink-deployer

VOLUME [ "/data/flink" ]

ENTRYPOINT [ "/flink-deployer/cli" ]
CMD [ "help" ]
