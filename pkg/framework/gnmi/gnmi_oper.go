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
	"time"

	"github.com/openconfig/gnmi/proto/gnmi"

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

//Init starts a gNMI client connection to switch under test
func Init(target *tg.Target) {
	log.Debug("In gnmi_oper Init")
	gnmiConn = connect(target)
	if gnmiConn.connError != nil {
		log.Fatalf("Unable to get a gnmi client: %v", gnmiConn.connError)
	}
}

//TearDown closes the gNMI connection
func TearDown() {
	log.Debug("In gnmi_oper TearDown")
	gnmiConn.cancel()
}

//ProcessGetRequest sends a request to switch and compares the response
func ProcessGetRequest(greq *gnmi.GetRequest, gresp *gnmi.GetResponse) bool {
	resp := gnmiConn.Get(greq)
	return verifyGetResp(gresp, resp)
}

//ProcessSetRequest sends a set request to switch and compares the response
func ProcessSetRequest(sreq *gnmi.SetRequest, sresp *gnmi.SetResponse) bool {
	resp := gnmiConn.Set(sreq)
	return verifySetResp(sresp, resp)
}

//ProcessSubscribeRequest opens a subscription channel to switch and processes the responses
func ProcessSubscribeRequest(sreq *gnmi.SubscribeRequest, sresp []*gnmi.SubscribeResponse, firstRespChan chan struct{}, resultChan chan bool) {
	subcl := gnmiConn.Subscribe()
	defer subcl.Close()
	log.Debugf("Length of expected result: %d\n\n", len(sresp))
	go subcl.Recv()
	go verifySubRespList(subcl.responseChan, sresp, firstRespChan, resultChan)
	if !subcl.Send(sreq) {
		resultChan <- false
	}

	select {
	case result := <-resultChan:
		resultChan <- result
	case <-time.After(SubTimeout):
		log.Error("Process subscribe request Timed out")
		resultChan <- false
	}

}

//verifySubResponses compares the responses from subscription channel with expected responses
func verifySubRespList(actRespChan chan *gnmi.SubscribeResponse, expResp []*gnmi.SubscribeResponse, firstRespChan chan struct{}, resultChan chan bool) {
	result, firstRespBool := true, true
	if len(expResp) == 0 {
		close(firstRespChan)
	} else {
		for _, exp := range expResp {
			act := <-actRespChan
			//Closing the first response channel to notify ProcessTelemetryExpectation to start processing actions
			if firstRespBool {
				firstRespBool = false
				close(firstRespChan)
			}
			log.Debug("In for loop after receiving subResp")
			result = verifySubResp(exp, act)
		}
	}
	resultChan <- result
}
