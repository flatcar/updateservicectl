FROM golang:1.6

WORKDIR /gopath/src/github.com/flatcar-linux/updateservicectl
ADD . /gopath/src/github.com/flatcar-linux/updateservicectl
RUN go get github.com/flatcar-linux/updateservicectl

CMD []
ENTRYPOINT ["/go/bin/updateservicectl"]
