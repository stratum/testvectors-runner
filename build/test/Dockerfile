# 
# Copyright 2019-present Open Networking Foundation
# 
# SPDX-License-Identifier: Apache-2.0
# 
# Start from golang v1.12.7 base image
FROM golang:1.12.7 as builder

RUN apt-get update \
  && apt-get install -y make libpcap-dev\
  && rm -rf /var/lib/apt/lists/* 

WORKDIR /root/testvectors-runner

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . . 

#Build go binary
RUN make build


# Docker image to run tvrunner with stratum_bmv2 switch
# Use "make switch" to start the switch and "make test" to run tests

FROM ubuntu:18.04

RUN apt-get update \
  && apt-get install -yq make libpcap-dev\
  && rm -rf /var/lib/apt/lists/* 

WORKDIR /root

RUN mkdir -p ./build/_output
COPY --from=builder /root/testvectors-runner/build/_output/tvrunner ./build/_output/tvrunner

RUN mkdir ./tools
COPY ./tools/Makefile ./tools/Makefile

RUN mkdir -p testvectors

WORKDIR /root

ENTRYPOINT ["/bin/bash"]

