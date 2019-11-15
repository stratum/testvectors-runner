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

//Connect description
func connect(tg *tvb.Target) connection {

	if tg.Address == "" {
		return connection{connError: errors.New("an address must be specified")}
	}
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, tg.Address, grpc.WithInsecure())
	if err != nil {
		return connection{connError: fmt.Errorf("cannot dial target %s, %v", tg.Address, err)}
	}

	return connection{ctx: ctx, client: gnmi.NewGNMIClient(conn), cancel: func() { conn.Close() }}
}
