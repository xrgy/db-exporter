FROM        quay.io/prometheus/busybox:latest
USER root


COPY db-exporter                          /bin/db-exporter
#ADD .                         $GOPATH/src/db-exporter
#WORKDIR $GOPATH/src/db-exporter


#RUN go build .
EXPOSE     9103
ENTRYPOINT ["/bin/db-exporter"]
