/*
 * Copyright 2019-present Open Networking Foundation
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package expectation

import (
	"time"

	"github.com/opennetworkinglab/testvectors-runner/pkg/framework/dataplane"
	"github.com/opennetworkinglab/testvectors-runner/pkg/framework/gnmi"
	"github.com/opennetworkinglab/testvectors-runner/pkg/framework/p4rt"
	"github.com/opennetworkinglab/testvectors-runner/pkg/logger"
	"github.com/opennetworkinglab/testvectors-runner/pkg/orchestrator/action"
	tv "github.com/stratum/testvectors/proto/testvector"
)

var log = logger.NewLogger()

//ProcessExpectation description
func ProcessExpectation(exp *tv.Expectation) bool {
	switch {
	case exp.GetConfigExpectation() != nil:
		ce := exp.GetConfigExpectation()
		return processConfigExpectation(ce)
	case exp.GetControlPlaneExpectation() != nil:
		cpe := exp.GetControlPlaneExpectation()
		return processControlPlaneExpectation(cpe)
	case exp.GetDataPlaneExpectation() != nil:
		dpe := exp.GetDataPlaneExpectation()
		return processDataPlaneExpectation(dpe)
	case exp.GetTelemetryExpectation() != nil:
		te := exp.GetTelemetryExpectation()
		return processTelemetryExpectation(te)
	default:
		log.Infof("Empty expectation\n")
		return false
	}
}

//processConfigExpectation extracts gnmi get requests and forwards it to framework
func processConfigExpectation(ce *tv.ConfigExpectation) bool {
	return gnmi.ProcessGetRequest(ce.GetGnmiGetRequest(), ce.GetGnmiGetResponse())
}

//processControlPlaneExpectation extracts get pipeline config, read or packet in expectations and forwards to framework.
func processControlPlaneExpectation(cpe *tv.ControlPlaneExpectation) bool {
	switch {
	case cpe.GetReadExpectation() != nil:
		log.Traceln("In Get Read Expectation")
		//TODO
	case cpe.GetPacketInExpectation() != nil:
		log.Traceln("In Get Packet In Expectation")
		return p4rt.ProcessPacketIn(cpe.GetPacketInExpectation().GetP4PacketIn())
	case cpe.GetPipelineConfigExpectation() != nil:
		log.Traceln("In Get Pipeline Config Expectation")
		//TODO
	}
	return false
}

//processDataPlaneExpectation extracts packets to be sent to data plane ports and forwards to framework.
func processDataPlaneExpectation(dpe *tv.DataPlaneExpectation) bool {
	switch {
	case dpe.GetTrafficExpectation() != nil:
		log.Traceln("In Get Traffic Expectation")
		// Get packet payloads
		pkts := dpe.GetTrafficExpectation().GetPackets()
		var payloads [][]byte
		for _, pkt := range pkts {
			payload := pkt.GetPayload()
			payloads = append(payloads, payload)
		}
		return dataplane.ProcessTrafficExpectation(payloads, dpe.GetTrafficExpectation().GetPorts())
	}
	return false
}

//processTelemetryExpectation executes subscribe expectations. These expectations contain gnmi subscribe request, set of actions to be performed after successful subscription and responses to be verfied.
//Returns false if responses are not received with in timeout.
func processTelemetryExpectation(tme *tv.TelemetryExpectation) bool {
	resultChan := make(chan bool, 1)
	firstRespChan := make(chan struct{})
	go gnmi.ProcessSubscribeRequest(tme.GetGnmiSubscribeRequest(), tme.GetGnmiSubscribeResponse(), firstRespChan, resultChan)
	select {
	case <-firstRespChan:
		actionResult := action.ProcessActionGroup(tme.GetActionGroup())
		select {
		case subResult := <-resultChan:
			log.Traceln("In ProcessTelemetryExpectation, Case Sub Result")
			return subResult && actionResult
		case <-time.After(gnmi.SubTimeout):
			log.Errorln("Timed out")
			return false
		}
	case <-time.After(gnmi.SubTimeout):
		log.Errorln("Timed out waiting for subscription response")
		return false
	}

}
