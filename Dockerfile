FROM golang:alpine as builder
LABEL maintainer="Thanabodee Charoenpiriyakij <wingyminus@gmail.com>"
RUN apk update && apk add git make
RUN mkdir -p $GOPATH/src/github.com/wingyplus/script_exporter
WORKDIR $GOPATH/src/github.com/wingyplus/script_exporter
COPY * $GOPATH/src/github.com/wingyplus/script_exporter/
RUN go get -v github.com/wingyplus/script_exporter/... \
    && go install -v github.com/wingyplus/script_exporter

FROM alpine
LABEL maintainer="Thanabodee Charoenpiriyakij <wingyminus@gmail.com>"
RUN mkdir -p /script_exporter/{bin,conf}
COPY --from=builder /go/bin/script_exporter /script_exporter/bin/
COPY --from=builder /go/src/github.com/wingyplus/script_exporter/scripts_config.yml /script_exporter/conf/
CMD ["/script_exporter/bin/script_exporter", "-scripts_config=/script_exporter/conf/scripts_config.yml"]
