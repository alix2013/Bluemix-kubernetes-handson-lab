FROM golang:1.8-alpine
COPY webserver.go /go/src/webserver/
RUN  apk add --no-cache git && go get -u github.com/go-redis/redis && go install webserver
ENV WEBAPP_PORT=8000 
EXPOSE $WEBAPP_PORT
CMD /go/bin/webserver

