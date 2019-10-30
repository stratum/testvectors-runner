# 
# Copyright 2019-present Open Networking Foundation
# 
# SPDX-License-Identifier: Apache-2.0
# 
PLATFORM := bmv2
TV_DIR := $$HOME/tv/${PLATFORM}/
TV_SUITE := p4runtime/
TV_NAME := RedirectDataplaneToDataplaneTest
TV_FILE := ${TV_DIR}${TV_SUITE}${TV_NAME}.pb.txt
TG_FILE := ${TV_DIR}target.pb.txt
PORT_MAP_FILE := tools/${PLATFORM}/port-map.json
LOG_LEVEL := fatal

switch:
	./tools/bmv2/start_bmv2.sh

pipeline:
	./tv_runner -test.v -tvFiles=${TV_DIR}PipelineConfig.pb.txt -tgFile=${TG_FILE} -portMapFile=${PORT_MAP_FILE} -logLevel=${LOG_LEVEL}

test:
	./tv_runner -test.v -tvFiles=${TV_FILE} -tgFile=${TG_FILE} -portMapFile=${PORT_MAP_FILE} -logLevel=${LOG_LEVEL}

gnmi:
	./tv_runner -test.v -tvDir=${TV_DIR}gnmi/ -tgFile=${TG_FILE} -portMapFile=${PORT_MAP_FILE} -logLevel=${LOG_LEVEL}

p4runtime:
	./tv_runner -test.v -tvDir=${TV_DIR}p4runtime/ -tgFile=${TG_FILE} -portMapFile=${PORT_MAP_FILE} -logLevel=${LOG_LEVEL}

e2e:
	./tv_runner -test.v -tvDir=${TV_DIR}e2e/ -tgFile=${TG_FILE} -portMapFile=${PORT_MAP_FILE} -logLevel=${LOG_LEVEL}

tests: pipeline gnmi e2e p4runtime
