/*
 * Copyright 2019-present Open Networking Foundation
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package gnmi

import (
	"time"

	"github.com/openconfig/gnmi/proto/gnmi"
)

//CtxTimeout for contexts
const CtxTimeout = 3 * time.Second

type subChan struct {
	client       gnmi.GNMI_SubscribeClient
	responseChan chan *gnmi.SubscribeResponse
}

//Close GNMI_SubscribeClient
func (s subChan) Close() {
	err := s.client.CloseSend()
	if err != nil {
		log.Warn("Error closing subscription client: ", err)
	}
}

//Recv runs a loop to continuously receive subscription responses from client and sends to specified channel.
//This method is called as go routine.
func (s subChan) Recv() {
	for {
		log.Debug("In Recv for loop")
		subResp, err := s.client.Recv()
		log.Debug("In Recv for loop after receiving message")
		if err != nil {
			log.Debugf("Failed to receive a message : %v\n", err)
			return
		}
		s.responseChan <- subResp
	}
}

//Send subscription requests to GNMI_SubscribeClient, returns false on error
func (s subChan) Send(subReq *gnmi.SubscribeRequest) bool {
	log.Info("Sending subscription request")
	log.Tracef("Subscription request: %s", subReq)
	err := s.client.Send(subReq)
	if err != nil {
		log.Errorf("Error sending subscription request: %v", err)
		return false
	}
	return true
}
