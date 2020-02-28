FROM alpine:latest
WORKDIR /root/
COPY ./dist/passman ./passman
RUN chmod +x ./passman
ENTRYPOINT [ "./main" ]
