FROM golang:1.12.1-alpine3.9
RUN mkdir -p /go/src/github.com/dbubel/passman
ADD . /go/src/github.com/dbubel/passman
WORKDIR /go/src/github.com/dbubel/passman/cmd/passman-api
RUN CGO_ENABLED=0 go build -a -v -ldflags '-extldflags "-static"' main.go
# RUN go build -v main.go
ENTRYPOINT [ "./main" ]