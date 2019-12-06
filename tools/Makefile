# 
# Copyright 2019-present Open Networking Foundation
# 
# SPDX-License-Identifier: Apache-2.0
# 
TVRUNNER_BIN := ../tvrunner
TV_DIR := $$HOME/testvectors/
TV_SUITE := p4runtime/
TV_NAME := RedirectDataplaneToDataplaneTest
TG_FILE := ${TV_DIR}target.pb.txt
PORT_MAP_FILE := ${TV_DIR}port-map.json
LOG_LEVEL := fatal
RUN_TVRUNNER := ${TVRUNNER_BIN} -test.v
ARGS := -target ${TG_FILE} -port-map ${PORT_MAP_FILE} -log-level ${LOG_LEVEL}

pipeline:
	${RUN_TVRUNNER} -tv-dir ${TV_DIR} -tv-name PipelineConfig ${ARGS}

test:
	${RUN_TVRUNNER} -tv-dir ${TV_DIR}${TV_SUITE} -tv-name ${TV_NAME} ${ARGS}

gnmi:
	${RUN_TVRUNNER} -tv-dir ${TV_DIR}gnmi/ ${ARGS}

p4runtime:
	${RUN_TVRUNNER} -tv-dir ${TV_DIR}p4runtime/ ${ARGS}

e2e:
	${RUN_TVRUNNER} -tv-dir ${TV_DIR}e2e/ ${ARGS}

tests: pipeline gnmi e2e p4runtime
