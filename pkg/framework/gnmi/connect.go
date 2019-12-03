/*
 * Copyright 2019-present Open Networking Foundation
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package gnmi

import (
	"context"
	"errors"
	"fmt"
	"time"

	"google.golang.org/grpc"

	"github.com/openconfig/gnmi/proto/gnmi"
	tvb "github.com/stratum/testvectors/proto/target"
)

//CtxTimeout for contexts
const CtxTimeout = 3 * time.Second

//connect starts a gRPC connection to the target specified.
//It returns connection struct with gNMI client, close function
//If an error is encountered during opening the connection, it is returned.
func connect(tg *tvb.Target) connection {

	if tg.Address == "" {
		return connection{connError: errors.New("an address must be specified")}
	}
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, CtxTimeout)
	defer cancel()
	conn, err := grpc.DialContext(ctx, tg.Address, grpc.WithInsecure())
	if err != nil {
		return connection{connError: fmt.Errorf("cannot dial target %s, %v", tg.Address, err)}
	}

	return connection{ctx: ctx, client: gnmi.NewGNMIClient(conn), cancel: func() { conn.Close() }}
}

//recvSubRespChan runs a loop to continuously receive subscription responses from client and sends to specified channel.
//This method is called as go routine.
func recvSubRespChan(subcl gnmi.GNMI_SubscribeClient, subRespChan chan *gnmi.SubscribeResponse) {
	for {
		log.Traceln("In recvSubRespChan for loop")
		subResp, err := subcl.Recv()
		log.Traceln("In recvSubRespChan for loop after receiving message")
		if err != nil {
			log.Tracef("Failed to receive a message : %v\n", err)
			return
		}
		subRespChan <- subResp
	}
}
