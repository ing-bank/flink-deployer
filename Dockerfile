FROM flink:1.3.2

ENV HIGH_AVAILABILITY=none \
    JOB_MANAGER_RPC_ADDRESS=localhost \
    JOB_MANAGER_RPC_PORT=6123 \
    JOB_MANAGER_WEB_ADDRESS=0.0.0.0 \
    JOB_MANAGER_WEB_PORT=8081

# Required to run our cli binary
RUN apt-get update && \
    apt-get install -y gettext vim musl

COPY ./conf/flink-conf-template.yaml /flink-deployer/conf/flink-conf-template.yaml
COPY docker-entrypoint.sh /flink-deployer/docker-entrypoint.sh
COPY cli /flink-deployer/cli

RUN mkdir -p /data/flink && \
    chgrp -R root /opt/flink && \
    chmod -R 775 /flink-deployer ${FLINK_HOME} /opt/flink

WORKDIR /flink-deployer

VOLUME [ "/data/flink" ]

ENTRYPOINT [ "/flink-deployer/docker-entrypoint.sh" ]
CMD [ "help" ]
