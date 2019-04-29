FROM golang:latest
RUN mkdir -p /go/src/github.com/dbubel/passman
ADD . /go/src/github.com/dbubel/passman
WORKDIR /go/src/github.com/dbubel/passman/cmd/passman-api
RUN go build -ldflags "\
                  -X main.GIT_HASH=`git rev-parse HEAD` \
                  -X main.BUILD_DATE=`date -u +'%Y-%m-%dT%H:%M:%SZ'`" \
                  -v main.go
# CGO_ENABLED=0 -a -v -ldflags '-extldflags "-static"'
# RUN go build -v main.go
ENTRYPOINT [ "./main" ]
