#!/usr/bin/env bash
# 
# Copyright 2019-present Open Networking Foundation
# 
# SPDX-License-Identifier: Apache-2.0
# 

NETWORK=host
IMAGE_NAME=stratumproject/tvrunner:binary

print_help() {
cat << EOF

Runs testvector based tests in a docker container with tvrunner binary. 
/tmp directory is mounted in the docker image to copy source files and logs.
tvrunner binary takes three mandatory arguments: target, port-map, tv-dir. 
Docker container starts in host network using default tvrunner:binary image.
The image name and network type can also be changed using additional arguments.

Usage: $0
    ***tvrunner arguments***
    [--target <filename>]               run testvectors against the provided target proto file
    [--port-map <filename>]             use the provided port mapping file
    [--tv-dir <directory>]              run all the testvectors from provided directory
    [--tv-name <regex>]                 run all the testvectors matching provided regular expression
    [--dp-mode <mode>]                  run the testvectors in provided mode
                                        default is direct; acceptable modes are <direct, loopbak>
    [--match-type <type>]               match packets based on the provided match-type
                                        default is exact; acceptable modes <exact, in>
    [--log-level <level>]               run tvrunner binary with provided log level
                                        default is error; acceptable levels are <panic, fatal, error, warn, info, debug, trace>
    [--log-dir <directory>]             save logs to provided directory
                                        default is /tmp

    ***docker arguments***
    [--image <name>]                    use the provided docker image
                                        default is $IMAGE_NAME
    [--network <name>]                  run tvrunner docker container in provided network
                                        default is $NETWORK

Examples:
    $0 --target ~/testvectors/bmv2/target.pb.txt --port-map ~/testvectors/bmv2/port-map.json --tv-dir ~/testvectors/bmv2/p4runtime
    $0 --target ~/testvectors/bmv2/target.pb.txt --port-map ~/testvectors/bmv2/port-map.json --tv-dir ~/testvectors/bmv2 --tv-name PipelineConfig
    $0 --target ~/testvectors/bmv2/target.pb.txt --port-map ~/testvectors/bmv2/port-map.json --tv-dir ~/testvectors/bmv2/p4runtime --tv-name PktIo.*
    $0 --target ~/testvectors/bmv2/target.pb.txt --port-map ~/testvectors/bmv2/port-map.json --tv-dir ~/testvectors/bmv2/p4runtime --image image:name --network none

EOF
}

while [[ $# -gt 0 ]]
do
    key="$1"
    case $key in
        -h|--help)
        print_help
        exit 0
        ;;
    --network)
        NETWORK="$2"
        shift 2
        ;;
    --image)
        IMAGE_NAME="$2"
        shift 2
        ;;
    --target)
        TG_FILE="$2"
        shift 2
        ;;
    --port-map)
        PM_FILE="$2"
        shift 2
        ;;
    --tv-dir)
        TV_DIR="$2"
        shift 2
        ;;
    --tv-name)
        TV_NAME="$2"
        shift 2
        ;;
    --dp-mode)
        DP_MODE="$2"
        shift 2
        ;;
    --match-type)
        MATCH_TYPE="$2"
        shift 2
        ;;
    --log-level)
        LOG_LEVEL="$2"
        shift 2
        ;;
    --log-dir)
        LOG_DIR="$2"
        shift 2
        ;;
    *)  # unknown option
        print_help
        exit 1
        ;;
    esac
done

# check mandatory arguments
if [[ -z $TG_FILE || -z $PM_FILE || -z $TV_DIR ]]; then
    print_help
    exit 1
fi

BINARY="build/_output/tvrunner "
TMP_TV_DIR=/tmp/tv_files

#make temporary directory
rm -rf $TMP_TV_DIR || true
mkdir -p $TMP_TV_DIR

#copy input files, folders
cp $TG_FILE $TMP_TV_DIR
cp $PM_FILE $TMP_TV_DIR
cp -rf $TV_DIR $TMP_TV_DIR

#
TG_FILE_MOUNT=$TMP_TV_DIR/${TG_FILE##*/}
PM_FILE_MOUNT=$TMP_TV_DIR/${PM_FILE##*/}
TV_DIR_MOUNT=$TMP_TV_DIR/${TV_DIR##*/}

DOCKER_RUN_OPTIONS="--rm -v /tmp:/tmp --network $NETWORK"
ENTRY_POINT="--entrypoint /root/$BINARY"

CMD="docker run $DOCKER_RUN_OPTIONS $ENTRY_POINT -ti $IMAGE_NAME"


TV_RUN_OPTIONS="-test.v --target $TG_FILE_MOUNT --port-map $PM_FILE_MOUNT --tv-dir $TV_DIR_MOUNT"

if [ -n "$TV_NAME" ]; then
    TV_RUN_OPTIONS="$TV_RUN_OPTIONS --tv-name $TV_NAME"
fi

if [ -n "$DP_MODE" ]; then
    TV_RUN_OPTIONS="$TV_RUN_OPTIONS --dp-mode $DP_MODE"
fi

if [ -n "$MATCH_TYPE" ]; then
    TV_RUN_OPTIONS="$TV_RUN_OPTIONS --match-type $MATCH_TYPE"
fi

if [ -n "$LOG_LEVEL" ]; then
    TV_RUN_OPTIONS="$TV_RUN_OPTIONS --log-level $LOG_LEVEL"
fi

if [ -n "$LOG_DIR" ]; then
    TV_RUN_OPTIONS="$TV_RUN_OPTIONS --log-dir $LOG_DIR"
fi

CMD="$CMD $TV_RUN_OPTIONS"

echo $CMD
$CMD

