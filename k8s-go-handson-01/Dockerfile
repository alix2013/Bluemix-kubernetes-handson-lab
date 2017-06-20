FROM golang:1.8-alpine
COPY webserver.go /go/src/webserver/
RUN go install webserver && rm -rf /go/src/webserver
ENV WEBAPP_PORT=8000 
EXPOSE $WEBAPP_PORT
CMD /go/bin/webserver

