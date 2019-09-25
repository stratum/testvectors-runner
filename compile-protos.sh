#!/bin/sh

proto_imports=".:${GOPATH}/src/github.com/googleapis/googleapis:${GOPATH}/src/github.com/abhilashendurthi/p4runtime/proto/:${GOPATH}/src"

protoc -I=$proto_imports --go_out=plugins=grpc:. pkg/proto/testvector/*.proto
protoc -I=$proto_imports --go_out=plugins=grpc:. pkg/proto/target/*.proto
