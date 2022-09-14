FROM golang:1.15

WORKDIR /gopath/src/github.com/flatcar/updateservicectl
ADD . /gopath/src/github.com/flatcar/updateservicectl
RUN go install github.com/flatcar/updateservicectl

CMD []
ENTRYPOINT ["/go/bin/updateservicectl"]
