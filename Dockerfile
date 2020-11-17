# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
ADD ./gpss /go/src/gpssclient/gpss
ADD . /go/src/kubernetes-postgres

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go get github.com/lib/pq
RUN go get github.com/golang/protobuf/proto
RUN go get google.golang.org/grpc
RUN go get github.com/streadway/amqp

WORKDIR /go/src/kubernetes-postgres
RUN go build 

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/src/gpss-rabbit-greenplum-connector/gpss-rabbit-connector

