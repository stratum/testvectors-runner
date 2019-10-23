/*
*Copyright 2019-present Open Networking Foundation
*
*SPDX-License-Identifier: Apache-2.0
 */

package orchestrator

import (
	"time"

	"github.com/opennetworkinglab/testvectors-runner/pkg/framework"
	tv "github.com/stratum/testvectors/proto/testvector"
)

//ProcessConfigExpectation will execute config expectations
func ProcessConfigExpectation(ce *tv.ConfigExpectation) bool {
	return framework.ProcessGetRequest(Target, ce.GetGnmiGetRequest(), ce.GetGnmiGetResponse())
}

//ProcessControlPlaneExpectation will execute control plane expectations
func ProcessControlPlaneExpectation(cpe *tv.ControlPlaneExpectation) bool {
	switch {
	case cpe.GetReadExpectation() != nil:
		log.Traceln("In Get Read Expectation")
		//TODO
	case cpe.GetPacketInExpectation() != nil:
		log.Traceln("In Get Packet In Expectation")
		return framework.ProcessPacketIn(cpe.GetPacketInExpectation().GetP4PacketIn())
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
		return framework.ProcessTrafficExpectation(Target, payloads, dpe.GetTrafficExpectation().GetPorts())
	}
	return false
}

//ProcessTelemetryExpectation will execute subscribe expectations
func ProcessTelemetryExpectation(tme *tv.TelemetryExpectation) bool {
	resultChan := make(chan bool, 1)
	var subResult, actionResult bool
	go framework.ProcessSubscribeRequest(Target, tme.GetGnmiSubscribeRequest(), tme.GetGnmiSubscribeResponse(), resultChan)
	time.Sleep(2 * time.Second)
	if ag := tme.GetActionGroup(); ag != nil {
		switch {
		case ag.GetSequentialActionGroup() != nil:
			actionResult = ProcessSequentialActionGroup(ag.GetSequentialActionGroup())
		case ag.GetParallelActionGroup() != nil:
			actionResult = ProcessParallelActionGroup(ag.GetParallelActionGroup())
		case ag.GetRandomizedActionGroup() != nil:
			actionResult = ProcessRandomizedActionGroup(ag.GetRandomizedActionGroup())
		default:
			log.Traceln("Empty Action Group")
			actionResult = false
		}
	}
	select {
	case subResult = <-resultChan:
		log.Traceln("In ProcessTelemetryExpectation, Case Sub Result")
		return subResult && actionResult
	case <-time.After(15 * time.Second):
		log.Errorln("Timed out")
		return false
	}
}
