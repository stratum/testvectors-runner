# Start from golang v1.12.7 base image
FROM golang:1.12.7

RUN apt-get update \
  && apt-get install -y vim libpcap-dev\
  && rm -rf /var/lib/apt/lists/*

RUN curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b $(go env GOPATH)/bin v1.21.0


# Set the Current Working Directory inside the container
WORKDIR /root/testvectors-runner

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .
