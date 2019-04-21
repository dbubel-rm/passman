FROM node:lts-alpine
RUN npm install -g http-server
RUN mkdir -p /trace-client
ADD . /trace-client
WORKDIR /trace-client

RUN npm run build

EXPOSE 8080
CMD [ "http-server", "dist" ]