FROM alpine:3.8

ENV FLINK_HOST=localhost \
    FLINK_PORT=6123

COPY cli /flink-deployer/cli

WORKDIR /flink-deployer

VOLUME [ "/data/flink" ]

ENTRYPOINT [ "/flink-deployer/cli" ]
CMD [ "help" ]
