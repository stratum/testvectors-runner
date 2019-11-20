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

//ProcessConfigExpectation will execute config expectations
func ProcessConfigExpectation(ce *tv.ConfigExpectation) bool {
	return gnmi.ProcessGetRequest(ce.GetGnmiGetRequest(), ce.GetGnmiGetResponse())
}

//ProcessControlPlaneExpectation will execute control plane expectations
func ProcessControlPlaneExpectation(cpe *tv.ControlPlaneExpectation) bool {
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

//ProcessDataPlaneExpectation will execute data plane expectations
func ProcessDataPlaneExpectation(dpe *tv.DataPlaneExpectation) bool {
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

//ProcessTelemetryExpectation will execute subscribe expectations
func ProcessTelemetryExpectation(tme *tv.TelemetryExpectation) bool {
	resultChan := make(chan bool, 1)
	var subResult, actionResult bool
	firstRespChan := make(chan struct{})
	go gnmi.ProcessSubscribeRequest(tme.GetGnmiSubscribeRequest(), tme.GetGnmiSubscribeResponse(), firstRespChan, resultChan)
	select {
	case <-firstRespChan:
		if ag := tme.GetActionGroup(); ag != nil {
			switch {
			case ag.GetSequentialActionGroup() != nil:
				actionResult = action.ProcessSequentialActionGroup(ag.GetSequentialActionGroup())
			case ag.GetParallelActionGroup() != nil:
				actionResult = action.ProcessParallelActionGroup(ag.GetParallelActionGroup())
			case ag.GetRandomizedActionGroup() != nil:
				actionResult = action.ProcessRandomizedActionGroup(ag.GetRandomizedActionGroup())
			default:
				log.Traceln("Empty Action Group")
				actionResult = false
			}
		}
		select {
		case subResult = <-resultChan:
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
