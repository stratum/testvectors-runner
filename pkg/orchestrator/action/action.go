/*
 * Copyright 2019-present Open Networking Foundation
 *
 * SPDX-License-Identifier: Apache-2.0
 */

/*
Package action implements functions to extract and run various actions supported by testvectors.
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

//ProcessActionGroup decodes the action group and executes actions sequentially, in parallel or randomly based on the type of underlying action group.
func ProcessActionGroup(ag *tv.ActionGroup) bool {
	log.Debug("In ProcessActionGroup")
	switch {
	case ag.GetSequentialActionGroup() != nil:
		sag := ag.GetSequentialActionGroup()
		return processSequentialActionGroup(sag)
	case ag.GetParallelActionGroup() != nil:
		pag := ag.GetParallelActionGroup()
		return processParallelActionGroup(pag)
	case ag.GetRandomizedActionGroup() != nil:
		rag := ag.GetRandomizedActionGroup()
		return processRandomizedActionGroup(rag)
	default:
		log.Info("Empty Action Group")
		return false
	}
}

//processSequentialActionGroup executes actions sequentially, combines all the results and returns a boolean value.
func processSequentialActionGroup(sag *tv.SequentialActionGroup) bool {
	result := true
	log.Debug("In ProcessSequentialActionGroup")
	for _, action := range sag.Actions {
		result = processAction(action) && result
	}
	return result
}

//processParallelActionGroup executes actions parallelly, combines all the results and returns a boolean value.
func processParallelActionGroup(pag *tv.ParallelActionGroup) bool {
	result := true
	log.Debug("In ProcessParallelActionGroup")
	//TODO - options
	//pag.Options
	var wg sync.WaitGroup
	wg.Add(len(pag.Actions))
	resultChan := make(chan bool, len(pag.Actions))
	for _, action := range pag.Actions {
		go func(action *tv.Action) {
			defer wg.Done()
			res := processAction(action)
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

//processRandomizedActionGroup executes actions in random order, combines all the results and returns a boolean value.
//TODO
func processRandomizedActionGroup(rag *tv.RandomizedActionGroup) bool {
	log.Debug("In ProcessRandomizedActionGroup")
	//TODO
	return false
}

//ProcessAction decodes and executes actions
func processAction(action *tv.Action) bool {
	log.Debug("In processAction")
	switch {
	case action.GetConfigOperation() != nil:
		co := action.GetConfigOperation()
		return processConfigOperation(co)
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
		log.Info("Empty Action")
	}
	return false
}

//processConfigOperation extracts gnmi set and forwards to framework.
func processConfigOperation(co *tv.ConfigOperation) bool {
	log.Debug("In processConfigOperation")
	return gnmi.ProcessSetRequest(co.GnmiSetRequest, co.GnmiSetResponse)
}

//processControlPlaneOperation extracts pipeline config, write or packet out operations and forwards to framework.
func processControlPlaneOperation(cpo *tv.ControlPlaneOperation) bool {
	log.Debug("In processControlPlaneOperation")
	switch {
	case cpo.GetPipelineConfigOperation() != nil:
		log.Debug("In Get Pipeline Config Oper")
		return p4rt.ProcessP4PipelineConfigOperation(cpo.GetPipelineConfigOperation().GetP4SetPipelineConfigRequest(), cpo.GetPipelineConfigOperation().GetP4SetPipelineConfigResponse())
	case cpo.GetWriteOperation() != nil:
		log.Debug("In Get Write Oper")
		return p4rt.ProcessP4WriteRequest(cpo.GetWriteOperation().GetP4WriteRequest(), cpo.GetWriteOperation().GetP4WriteResponse())
	case cpo.GetPacketOutOperation() != nil:
		log.Debug("In PacketOut Oper")
		return p4rt.ProcessPacketOutOperation(cpo.GetPacketOutOperation().GetP4PacketOut())
	}
	return false
}

//ProcessDataPlaneStimulus extracts traffic stimulus and forwards to framework.
func processDataPlaneStimulus(dps *tv.DataPlaneStimulus) bool {
	log.Debug("in processDataPlaneStimulus")
	switch {
	case dps.GetTrafficStimulus() != nil:
		log.Debug("In Get Traffic Stimulus")
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
