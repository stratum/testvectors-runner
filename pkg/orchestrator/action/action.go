/*
 * Copyright 2019-present Open Networking Foundation
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package action

import (
	"sync"

	"github.com/opennetworkinglab/testvectors-runner/pkg/framework/dataplane"
	"github.com/opennetworkinglab/testvectors-runner/pkg/framework/gnmi"
	"github.com/opennetworkinglab/testvectors-runner/pkg/framework/p4rt"
	"github.com/opennetworkinglab/testvectors-runner/pkg/logger"
	tv "github.com/stratum/testvectors/proto/testvector"
)

var log = logger.NewLogger()

//ProcessSequentialActionGroup will execute actions sequentially
func ProcessSequentialActionGroup(sag *tv.SequentialActionGroup) bool {
	result := true
	log.Traceln("In ProcessSequentialActionGroup")
	for _, action := range sag.Actions {
		result = ProcessAction(action) && result
	}
	return result
}

//ProcessParallelActionGroup will execute actions Parallelly
func ProcessParallelActionGroup(pag *tv.ParallelActionGroup) bool {
	result := true
	log.Traceln("In ProcessParallelActionGroup")
	//TODO - options
	//pag.Options
	var wg sync.WaitGroup
	wg.Add(len(pag.Actions))
	resultChan := make(chan bool, len(pag.Actions))
	for _, action := range pag.Actions {
		go func(action *tv.Action) {
			defer wg.Done()
			res := ProcessAction(action)
			resultChan <- res
		}(action)
	}
	wg.Wait()
	close(resultChan)
	for r := range resultChan {
		result = result && r
	}
	return result
}

//ProcessRandomizedActionGroup will execute actions in random order
func ProcessRandomizedActionGroup(rag *tv.RandomizedActionGroup) bool {
	log.Traceln("In ProcessRandomizedActionGroup")
	//TODO
	return false
}

//ProcessAction decodes and executes actions
func ProcessAction(action *tv.Action) bool {
	switch {
	case action.GetConfigOperation() != nil:
		co := action.GetConfigOperation()
		return gnmi.ProcessSetRequest(co.GnmiSetRequest, co.GnmiSetResponse)
	case action.GetAlarmStimulus() != nil:
		//TODO
		as := action.GetAlarmStimulus()
		_ = as
	case action.GetControlPlaneOperation() != nil:
		//TODO
		cpo := action.GetControlPlaneOperation()
		return processControlPlaneOperation(cpo)
	case action.GetDataPlaneStimulus() != nil:
		//TODO
		dps := action.GetDataPlaneStimulus()
		return processDataPlaneStimulus(dps)
	case action.GetManagementOperation() != nil:
		//TODO
		mo := action.GetManagementOperation()
		_ = mo
	case action.GetPortStimulus() != nil:
		//TODO
		ps := action.GetPortStimulus()
		_ = ps
	default:
		log.Traceln("Empty Action")
	}
	return false
}

//ProcessControlPlaneOperation extracts pipeline config, write or packet out operations and forwards to framework.
func processControlPlaneOperation(cpo *tv.ControlPlaneOperation) bool {
	switch {
	case cpo.GetPipelineConfigOperation() != nil:
		log.Traceln("In Get Pipeline Config Oper")
		return p4rt.ProcessP4PipelineConfigOperation(cpo.GetPipelineConfigOperation().GetP4SetPipelineConfigRequest(), cpo.GetPipelineConfigOperation().GetP4SetPipelineConfigResponse())
	case cpo.GetWriteOperation() != nil:
		log.Traceln("In Get Write Oper")
		return p4rt.ProcessP4WriteRequest(cpo.GetWriteOperation().GetP4WriteRequest(), cpo.GetWriteOperation().GetP4WriteResponse())
	case cpo.GetPacketOutOperation() != nil:
		log.Traceln("In PacketOut Oper")
		return p4rt.ProcessPacketOutOperation(cpo.GetPacketOutOperation().GetP4PacketOut())
	}
	return false
}

//ProcessDataPlaneStimulus extracts traffic stimulus and forwards to framework.
func processDataPlaneStimulus(dps *tv.DataPlaneStimulus) bool {
	switch {
	case dps.GetTrafficStimulus() != nil:
		log.Traceln("In Get Traffic Stimulus")
		// Get packet payloads
		pkts := dps.GetTrafficStimulus().GetPackets()
		var payloads [][]byte
		for _, pkt := range pkts {
			payload := pkt.GetPayload()
			payloads = append(payloads, payload)
		}
		return dataplane.ProcessTrafficStimulus(payloads, dps.GetTrafficStimulus().GetPort())
	}
	return false
}
