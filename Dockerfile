FROM golang:1.15

WORKDIR /gopath/src/github.com/kinvolk/updateservicectl
ADD . /gopath/src/github.com/kinvolk/updateservicectl
RUN go install github.com/kinvolk/updateservicectl

CMD []
ENTRYPOINT ["/go/bin/updateservicectl"]
