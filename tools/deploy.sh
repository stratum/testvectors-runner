#!/bin/bash
# 
# Copyright 2019-present Open Networking Foundation
# 
# SPDX-License-Identifier: Apache-2.0
# 

REMOTE_TVRUNNER_DIR="~/testvectors-runner"
TVRUNNER_BIN=../build/_output/tvrunner

# Create tvrunner directory on remote node
ssh -tt $1 "
	[ -d $REMOTE_TVRUNNER_DIR ] || mkdir $REMOTE_TVRUNNER_DIR
	[ -d $REMOTE_TVRUNNER_DIR/build/_output ] || mkdir -p $REMOTE_TVRUNNER_DIR/build/_output
	[ -d $REMOTE_TVRUNNER_DIR/tools ] || mkdir -p $REMOTE_TVRUNNER_DIR/tools
"

# Copy tvrunner binary and other files to remote node
scp $TVRUNNER_BIN $1:$REMOTE_TVRUNNER_DIR/build/_output/tvrunner
scp ./Makefile $1:$REMOTE_TVRUNNER_DIR/tools
