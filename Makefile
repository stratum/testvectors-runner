# 
# Copyright 2019-present Open Networking Foundation
# 
# SPDX-License-Identifier: Apache-2.0
# 
TV_DIR := $$HOME/testvectors/bmv2
TV_RUNNER_DIR := $$PWD

.PHONY: build

build:
	CGO_ENABLED=1 go build -o build/_output/tv_runner ./cmd/main

bmv2:
	docker run --privileged --rm -it -p50001:50001 --name bmv2  stratumproject/tvrunner:bmv2

tv-runner-dev: #WIP
	docker run --rm -it --network=container:bmv2 -v ${TV_RUNNER_DIR}:/root/testvectors-runner -v ${TV_DIR}:/root/tv/bmv2 stratumproject/tvrunner:dev

tv-runner:
	docker run --rm -it --network=container:bmv2 -v ${TV_DIR}:/root/tv/bmv2 stratumproject/tvrunner:binary

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

test: build deps linters license_check
	CGO_ENABLED=1 go test -race -v github.com/opennetworkinglab/testvectors-runner/pkg/framework/...
	CGO_ENABLED=1 go test -race -v github.com/opennetworkinglab/testvectors-runner/pkg/orchestrator/...
