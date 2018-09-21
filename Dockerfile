FROM golang:1.10.4-alpine3.8
WORKDIR /go/src/github.com/ing-bank/flink-deployer/
COPY . .
RUN go build ./cmd/cli

FROM alpine:3.8
WORKDIR /flink-deployer
COPY --from=0 /go/src/github.com/ing-bank/flink-deployer/cli .
VOLUME [ "/data/flink" ]
ENTRYPOINT [ "/flink-deployer/cli" ]
CMD [ "help" ]
