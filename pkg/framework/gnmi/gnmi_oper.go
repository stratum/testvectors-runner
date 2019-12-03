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
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/openconfig/ygot/testutil"

	"github.com/opennetworkinglab/testvectors-runner/pkg/logger"
	tg "github.com/stratum/testvectors/proto/target"
)

var (
	log      = logger.NewLogger()
	gnmiConn connection
)

const (
	//SubTimeout for receiving subscription acknowledgement
	SubTimeout = 5 * time.Second
)

//Connection struct stores the gNMI client connection, context and cancel function.
type connection struct {
	ctx       context.Context
	client    gnmi.GNMIClient
	connError error
	cancel    context.CancelFunc
}

//Init starts a gNMI client connection to switch under test
func Init(target *tg.Target) {
	gnmiConn = connect(target)
	if gnmiConn.connError != nil {
		log.Errorln(gnmiConn.connError)
		log.Fatalln("Unable to get a gnmi client")
	}
}

//TearDown closes the gNMI connection
func TearDown() {
	log.Traceln("In gnmi_oper tear down")
	gnmiConn.cancel()
}

//ProcessGetRequest sends a request to switch and compares the response
func ProcessGetRequest(greq *gnmi.GetRequest, gresp *gnmi.GetResponse) bool {
	log.Infoln("Sending get request")
	ctx := context.Background()
	resp, err := gnmiConn.client.Get(ctx, greq)
	if err != nil {
		log.Errorln(err)
		return false
	}
	isEqual := testutil.GetResponseEqual(resp, gresp, testutil.IgnoreTimestamp{})
	if !isEqual {
		log.Warningf("Get responses are unequal\nExpected: %s\nActual  : %s\n", gresp, resp)
	} else {
		log.Infoln("Get responses are equal")
		log.Debugf("Get response: %s\n", resp)
	}
	return isEqual
}

//ProcessSetRequest sends a set request to switch and compares the response
func ProcessSetRequest(sreq *gnmi.SetRequest, sresp *gnmi.SetResponse) bool {
	log.Traceln("In ProcessSetRequest")
	log.Infof("Sending set request: %s", sreq)
	ctx := context.Background()

	resp, err := gnmiConn.client.Set(ctx, sreq)
	if err != nil {
		log.Errorln(err)
		return false
	}
	//FIXME
	//resetting timestamp as a work around to ignore timestamp during comparison
	origRespTimestamp := resp.Timestamp
	origSRespTimestamp := sresp.Timestamp
	resp.Timestamp = 0
	sresp.Timestamp = 0
	isEqual := proto.Equal(resp, sresp)
	if !isEqual {
		log.Warningf("Set responses are unequal\nexpected: %s\nactual: %s\n", sresp, resp)
	} else {
		log.Infoln("Set responses are equal")
		log.Debugf("Set response: %s\n", resp)
	}

	//Reset timestamp to original value
	resp.Timestamp = origRespTimestamp
	sresp.Timestamp = origSRespTimestamp
	return isEqual
}

//ProcessSubscribeRequest opens a subscription channel to switch and processes the responses
func ProcessSubscribeRequest(sreq *gnmi.SubscribeRequest, sresp []*gnmi.SubscribeResponse, firstRespChan chan struct{}, resultChan chan bool) {
	ctx := context.Background()
	subcl, err := gnmiConn.client.Subscribe(ctx)
	if err != nil {
		log.Infoln(err)
	}
	defer func() {
		err := subcl.CloseSend()
		if err != nil {
			log.Warnln("Error closing subscription client: ", err)
		}
	}()
	log.Tracef("Length of expected result: %d\n\n", len(sresp))
	var result bool
	subRespChan := make(chan *gnmi.SubscribeResponse)
	go recvSubRespChan(subcl, subRespChan)
	go verifySubResponses(subRespChan, sresp, firstRespChan, resultChan)
	log.Infoln("Sending subscription request")
	err = subcl.Send(sreq)

	if err != nil {
		log.Errorln(err)
		resultChan <- false
	}

	select {
	case result = <-resultChan:
		resultChan <- result
	case <-time.After(SubTimeout):
		log.Errorln("Process subscribe request Timed out")
		resultChan <- false
	}

}

//verifySubResponses compares the responses from subscription channel with expected responses
func verifySubResponses(actRespChan chan *gnmi.SubscribeResponse, expResp []*gnmi.SubscribeResponse, firstRespChan chan struct{}, resultChan chan bool) {
	result, firstRespBool := true, true
	if len(expResp) == 0 {
		close(firstRespChan)
	} else {
		for _, exp := range expResp {
			log.Traceln("In for loop")
			act := <-actRespChan
			if firstRespBool {
				firstRespBool = false
				close(firstRespChan)
			}
			log.Traceln("In for loop after receiving subResp")
			if act.GetUpdate() != nil && testutil.NotificationSetEqual([]*gnmi.Notification{act.GetUpdate()}, []*gnmi.Notification{exp.GetUpdate()}, testutil.IgnoreTimestamp{}) {
				log.Traceln("In GetUpdate condition")
				log.Infoln("Subscription responses are equal")
				log.Debugf("Subscription response: %s\n", act)
				result = result && true
			} else if testutil.SubscribeResponseEqual(exp, act) {
				//continue
				log.Infoln("Subscription responses are equal")
				log.Debugf("Subscription response: %s\n", act)
				result = result && true
			} else {
				log.Warnf("Subscription responses are unequal expected:\n%s \nactual:\n%s", exp, act)
				result = result && false
			}
		}
	}
	resultChan <- result
}
