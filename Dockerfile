FROM golang:1.11.5-alpine3.9 as build
WORKDIR /go/src/github.com/ing-bank/flink-deployer/
COPY . .
RUN go build ./cmd/cli

FROM alpine:3.9
WORKDIR /flink-deployer
COPY --from=build /go/src/github.com/ing-bank/flink-deployer/cli .
VOLUME [ "/data/flink" ]
ENTRYPOINT [ "/flink-deployer/cli" ]
CMD [ "help" ]
