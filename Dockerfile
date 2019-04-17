FROM golang:1.12.1-alpine3.9
RUN mkdir -p /go/src/github.com/dbubel/passman
ADD . /go/src/github.com/dbubel/passman
WORKDIR /go/src/github.com/dbubel/passman/cmd/passman-api
RUN go build main.go
# CGO_ENABLED=0 -a -v -ldflags '-extldflags "-static"' 
# RUN go build -v main.go
ENTRYPOINT [ "./main" ]