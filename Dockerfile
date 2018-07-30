FROM alpine:3.8

COPY cli /flink-deployer/cli

WORKDIR /flink-deployer

VOLUME [ "/data/flink" ]

ENTRYPOINT [ "/flink-deployer/cli" ]
CMD [ "help" ]
