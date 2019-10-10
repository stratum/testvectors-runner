# Start from golang v1.12.7 base image
FROM golang:1.12.7 as builder

RUN apt-get update \
  && apt-get install -y vim libpcap-dev\
  && rm -rf /var/lib/apt/lists/*

#Work around to access testvectors(https://github.com/opennetworkinglab/testvectors) private git repository
ARG GIT_USER
ARG GIT_PERSONAL_ACCESS_TOKEN
RUN git config \
    --global \
    url."https://${GIT_USER}:${GIT_PERSONAL_ACCESS_TOKEN}@github.com/opennetworkinglab/testvectors".insteadOf \
    "https://github.com/opennetworkinglab/testvectors"
# Set the Current Working Directory inside the container
WORKDIR /root/testvectors-runner

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . . 

#Build go binary
RUN env GOOS=$(go env GOOS) GOARCH=$(go env GOARCH) go build -o ./cmd/main/tv_runner ./cmd/main/testvectors-runner.go



# Docker image to run tv_runner with stratum_bmv2 switch
# Use "make switch" to start the switch and "make test" to run tests

FROM opennetworking/mn-stratum

RUN apt-get update && apt-get install -yq make vim

WORKDIR /root

COPY --from=builder /root/testvectors-runner/cmd/main/tv_runner ./tv_runner
COPY ./tools/Makefile ./Makefile

RUN mkdir ./tools
COPY ./tools/bmv2 ./tools/bmv2/

RUN mkdir -p tv/bmv2

WORKDIR /root

ENTRYPOINT ["/bin/bash"]

