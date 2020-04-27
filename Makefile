# 
# Copyright 2019-present Open Networking Foundation
# 
# SPDX-License-Identifier: Apache-2.0
# 
TV_DIR := $$HOME/testvectors/bmv2
TVRUNNER_DIR := $$PWD
DOCKER_TV_DIR := /root/testvectors
DOCKER_TVRUNNER_DIR := /root/testvectors-runner
DOCKER_RUN := docker run --rm -it
DOCKER_RUN_BMV2 := ${DOCKER_RUN} --network=container:bmv2
DOCKER_RUN_HW := ${DOCKER_RUN} --network=host
TVRUNNER_DEV_IMAGE := stratumproject/tvrunner:dev
TVRUNNER_BIN_IMAGE := stratumproject/tvrunner:binary

.PHONY: build

build:
	CGO_ENABLED=1 go build -o tvrunner ./cmd/main

bmv2:
	${DOCKER_RUN} --privileged -p50001:50001 --name bmv2 --network=host stratumproject/tvrunner:bmv2

tvrunner-bmv2-dev:
	${DOCKER_RUN_BMV2} -v ${TVRUNNER_DIR}:${DOCKER_TVRUNNER_DIR} -v ${TV_DIR}:${DOCKER_TV_DIR} ${TVRUNNER_DEV_IMAGE}

tvrunner-hw-dev:
	${DOCKER_RUN_HW} -v ${TVRUNNER_DIR}:${DOCKER_TVRUNNER_DIR} -v ${TV_DIR}:${DOCKER_TV_DIR} ${TVRUNNER_DEV_IMAGE}

tvrunner-bmv2:
	${DOCKER_RUN_BMV2} -v ${TV_DIR}:${DOCKER_TV_DIR} ${TVRUNNER_BIN_IMAGE}

tvrunner-hw:
	${DOCKER_RUN_HW} -v ${TV_DIR}:${DOCKER_TV_DIR} ${TVRUNNER_BIN_IMAGE}

deps: # @HELP ensure that the required dependencies are in place
	go build -v ./...
	bash -c "diff -u <(echo -n) <(git diff go.mod)"
	bash -c "diff -u <(echo -n) <(git diff go.sum)"

linters: # @HELP examines Go source code and reports coding problems
	golangci-lint run

gofmt: # @HELP run the Go format validation
	bash -c "diff -u <(echo -n) <(gofmt -d pkg/ cmd/ tests/)"

license_check: # @HELP examine and ensure license headers exist
	./build/licensing/boilerplate.py -v

test: build linters license_check
	CGO_ENABLED=1 go test -race -v github.com/stratum/testvectors-runner/pkg/framework/...
	CGO_ENABLED=1 go test -race -v github.com/stratum/testvectors-runner/pkg/orchestrator/...
