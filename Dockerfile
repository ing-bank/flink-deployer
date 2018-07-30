FROM alpine:3.8

ENV JOB_MANAGER_RPC_ADDRESS=localhost \
    JOB_MANAGER_RPC_PORT=6123

COPY cli /flink-deployer/cli

WORKDIR /flink-deployer

VOLUME [ "/data/flink" ]

ENTRYPOINT [ "/flink-deployer/cli" ]
CMD [ "help" ]
