FROM golang:latest
RUN mkdir -p /go/src/github.com/dbubel/passman
ADD . /go/src/github.com/dbubel/passman
WORKDIR /go/src/github.com/dbubel/passman/cmd/passman-api
#RUN CGO_ENABLED=0 go build -a -v -ldflags '-extldflags "-static"' main.go
ENV PORT=80
RUN go build -v main.go
ENTRYPOINT [ "./main" ]