FROM golang:1.17 as builder

COPY . /demo

WORKDIR /demo
RUN go build -o bin/prometheus-minimal-demo main.go

FROM debian:stretch

COPY --from=builder \
  /demo/bin/prometheus-minimal-demo \
  /usr/local/bin/prometheus-minimal-demo

USER        nobody
ENTRYPOINT  ["prometheus-minimal-demo"]
