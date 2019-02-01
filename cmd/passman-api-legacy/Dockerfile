FROM golang:alpine
RUN mkdir -p /go/src/github.com/dbubel/passman
ADD . /go/src/github.com/dbubel/passman
WORKDIR /go/src/github.com/dbubel/passman
#RUN CGO_ENABLED=0 go build -a -v -ldflags '-extldflags "-static"' main.go
ENV PORT=80
RUN go build -v main.go
ENTRYPOINT [ "./main" ]
# $(aws ecr get-login --no-include-email --region us-east-1)
# docker build -t passman-production-ecr .
# docker tag passman-production-ecr:latest 316188497159.dkr.ecr.us-east-1.amazonaws.com/passman-production-ecr:latest
# docker push 316188497159.dkr.ecr.us-east-1.amazonaws.com/passman-production-ecr:latest
