/*
 * Copyright 2019-present Open Networking Foundation
 *
 * SPDX-License-Identifier: Apache-2.0
 */

/*
Package gnmi implements gnmi get, set and subscribe functions
*/
package gnmi

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/openconfig/ygot/testutil"
	tvb "github.com/stratum/testvectors/proto/target"
	"google.golang.org/grpc"
)

//Connection struct stores the gNMI client connection, context and cancel function.
type connection struct {
	ctx       context.Context
	client    gnmi.GNMIClient
	connError error
	cancel    context.CancelFunc
}

//connect starts a gRPC connection to the target specified.
//It returns connection struct with gNMI client, close function
//If an error is encountered during opening the connection, it is returned.
func connect(tg *tvb.Target) connection {
	log.Debug("In gnmi_oper connect")
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

//Get calls gNMI client's Get RPC call and returns the GetResponse
func (c *connection) Get(getReq *gnmi.GetRequest) *gnmi.GetResponse {
	log.Info("Sending get request")
	ctx := context.Background()
	resp, err := c.client.Get(ctx, getReq)
	if err != nil {
		log.Error(err)
		return nil
	}
	return resp
}

//Get calls gNMI client's Set RPC call and returns the SetResponse
func (c *connection) Set(setReq *gnmi.SetRequest) *gnmi.SetResponse {
	log.Info("Sending set request")
	log.Debugf("Set request: %s", setReq)
	ctx := context.Background()

	resp, err := c.client.Set(ctx, setReq)
	if err != nil {
		log.Errorf("Error sending set request: %v", err)
		return nil
	}
	return resp
}

//Subscribe calls gNMI client's Subscribe RPC call and returns struct with GNMI_SubscribeClient and SubscribeResponse channel
func (c *connection) Subscribe() subChan {
	ctx := context.Background()
	client, err := c.client.Subscribe(ctx)
	if err != nil {
		log.Errorf("Error getting subscription client: %v", err)
		return subChan{}
	}
	return subChan{client: client, responseChan: make(chan *gnmi.SubscribeResponse)}
}

//verifyGetResp compares two gnmi GetResponses and returns true or false
func verifyGetResp(expected, actual *gnmi.GetResponse) bool {
	switch {
	case expected == nil && actual == nil:
		log.Debug("Both get responses are empty")
		return true
	case expected == nil || actual == nil:
		log.Warnf("Get responses are unequal\nExpected: %s\nActual  : %s\n", expected, actual)
		return false
	case testutil.GetResponseEqual(expected, actual, testutil.IgnoreTimestamp{}):
		log.Info("Get responses are equal")
		log.Debugf("Get response: %s\n", actual)
		return true
	default:
		log.Warnf("Get responses are unequal\nExpected: %s\nActual  : %s\n", expected, actual)
		return false
	}
}

//verifySetResp compares two gnmi SetResponses and returns true or false
func verifySetResp(expected, actual *gnmi.SetResponse) bool {
	switch {
	case expected == nil && actual == nil:
		log.Debug("Both set responses are empty")
		return true
	case expected == nil || actual == nil:
		log.Warnf("Set responses are unequal\nExpected: %s\nActual  : %s\n", expected, actual)
		return false
	default:
		//FIXME
		//resetting timestamp as a work around to ignore timestamp during comparison
		expected.Timestamp = 0
		actual.Timestamp = 0
		if proto.Equal(expected, actual) {
			log.Info("Set responses are equal")
			log.Debugf("Set response: %s", actual)
			return true
		}
		log.Warnf("Set responses are unequal\nexpected: %s\nactual: %s\n", expected, actual)
		return false
	}
}

//verifySubResp compares two gnmi SubscribeResponses and returns true or false
func verifySubResp(expected, actual *gnmi.SubscribeResponse) bool {
	switch {
	case expected == nil && actual == nil:
		log.Debug("Both set responses are empty")
		return true
	case expected == nil || actual == nil:
		log.Warnf("Set responses are unequal\nExpected: %s\nActual  : %s\n", expected, actual)
		return false
	case actual.GetUpdate() != nil && testutil.NotificationSetEqual([]*gnmi.Notification{expected.GetUpdate()}, []*gnmi.Notification{actual.GetUpdate()}, testutil.IgnoreTimestamp{}):
		log.Debug("In GetUpdate condition")
		log.Info("Subscription responses are equal")
		log.Debugf("Subscription response: %s\n", actual)
		return true
	case testutil.SubscribeResponseEqual(expected, actual):
		//continue
		log.Info("Subscription responses are equal")
		log.Debugf("Subscription response: %s\n", actual)
		return true
	default:
		log.Warnf("Subscription responses are unequal expected:\n%s \nactual:\n%s", expected, actual)
		return false
	}

}
