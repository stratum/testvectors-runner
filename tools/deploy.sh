#!/bin/bash

PLATFORM=tofino
REMOTE_TV_RUNNER_DIR="~/tv_runner"
TV_RUNNER_BIN=../cmd/main/tv_runner

# Create tv_runner directory on remote node
ssh -tt $1 "
	[ -d $REMOTE_TV_RUNNER_DIR ] || mkdir $REMOTE_TV_RUNNER_DIR
	[ -d $REMOTE_TV_RUNNER_DIR/tools/$PLATFORM ] || mkdir -p $REMOTE_TV_RUNNER_DIR/tools/$PLATFORM
"

# Copy tv_runner binary and other files to remote node
scp $TV_RUNNER_BIN $1:$REMOTE_TV_RUNNER_DIR/tv_runner
scp ./Makefile $1:$REMOTE_TV_RUNNER_DIR
scp ./$PLATFORM/port-map.json $1:$REMOTE_TV_RUNNER_DIR/tools/$PLATFORM

# Change platform name in Makefile
ssh -tt $1 "sed -i '/PLATFORM :=/c\PLATFORM := $PLATFORM' $REMOTE_TV_RUNNER_DIR/Makefile"

