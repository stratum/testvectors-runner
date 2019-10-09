/*
Copyright 2018 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package common defines operations that are used within the testvectors-runner
// framework for multiple tests.
// Copied from github.com/openconfig/gnmitest/common/gnmi.go
package common

import (
	"context"
	"errors"
	"fmt"
	"time"

	"google.golang.org/grpc"

	v1 "github.com/abhilashendurthi/p4runtime/proto/p4/v1"
	gpb "github.com/openconfig/gnmi/proto/gnmi"

	"github.com/opennetworkinglab/testvectors-runner/pkg/logger"
	tvb "github.com/opennetworkinglab/testvectors/proto/target"
)

var log = logger.NewLogger()

// Connect opens a new gRPC connection to the target speciifed by the
// ConnectionArgs. It returns the gNMI Client connection, and a function
// which can be called to close the connection. If an error is encountered
// during opening the connection, it is returned.
func Connect(ctx context.Context, a *tvb.Target) (gpb.GNMIClient, func(), error) {
	if a.Address == "" {
		return nil, nil, errors.New("an address must be specified")
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	/*conn, err := grpc.DialContext(ctx, a.Address, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: true,
	})))*/
	conn, err := grpc.DialContext(ctx, a.Address, grpc.WithInsecure())
	if err != nil {
		return nil, nil, fmt.Errorf("cannot dial target %s, %v", a.Address, err)
	}

	return gpb.NewGNMIClient(conn), func() { conn.Close() }, nil
}

// ConnectP4 opens a new gRPC connection to the target speciifed by the
// ConnectionArgs. It returns the p4runtime Client connection, and a function
// which can be called to close the connection. If an error is encountered
// during opening the connection, it is returned.
func ConnectP4(ctx context.Context, a *tvb.Target) (v1.P4RuntimeClient, func(), error) {
	if a.Address == "" {
		return nil, nil, errors.New("an address must be specified")
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, a.Address, grpc.WithInsecure())
	if err != nil {
		return nil, nil, fmt.Errorf("cannot dial target %s, %v", a.Address, err)
	}

	return v1.NewP4RuntimeClient(conn), func() { conn.Close() }, nil

}
