# Start from golang v1.12.7 base image
FROM golang:1.12.7

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
#COPY . .
