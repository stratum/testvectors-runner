# Start from golang v1.12.7 base image
FROM golang:1.12.7 as builder

RUN apt-get update \
  && apt-get install -y vim libpcap-dev\
  && rm -rf /var/lib/apt/lists/* 

WORKDIR /root/testvectors-runner

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . . 

#Build go binary
RUN env GOOS=$(go env GOOS) GOARCH=$(go env GOARCH) go build -o ./cmd/main/tv_runner ./cmd/main/testvectors-runner.go



# Docker image to run tv_runner with stratum_bmv2 switch
# Use "make switch" to start the switch and "make test" to run tests

FROM ubuntu:18.04

RUN apt-get update \
  && apt-get install -yq make libpcap-dev\
  && rm -rf /var/lib/apt/lists/* 

WORKDIR /root

COPY --from=builder /root/testvectors-runner/cmd/main/tv_runner ./tv_runner
COPY ./tools/Makefile ./Makefile

RUN mkdir ./tools
COPY ./tools/bmv2 ./tools/bmv2/

RUN mkdir -p tv/bmv2

WORKDIR /root

ENTRYPOINT ["/bin/bash"]

