# Start from golang v1.12 base image
FROM golang:1.12 as builder

# Add Maintainer Info
LABEL maintainer="Ilya Dokshukin <dokshukin@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR /go/src/github.com/dokshukin/gl2j-gw
RUN apt-get update && \
    apt-get install -y unzip upx && \
    wget https://github.com/dokshukin/gl2j-gw/archive/master.zip && \
    unzip master.zip -d /tmp/ && \
    mv /tmp/gl2j-gw-master/* .  && \
    rm master.zip

# Download dependencies
RUN go get -d -v ./...

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build --ldflags "-s -w -X main.version=${DOCKER_TAG} -X main.build=2019-06-04" -a -installsuffix cgo -o /go/bin/gl2j-gw .
RUN upx /go/bin/gl2j-gw

######## Start a new stage from scratch #######
FROM alpine:latest

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /go/bin/gl2j-gw .

EXPOSE 8080

CMD ["./gl2j-gw", "-c", "/etc/config.yml"]