/*
 * Copyright 2019-present Open Networking Foundation
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package testvector

import (
	"github.com/opennetworkinglab/testvectors-runner/pkg/logger"
	"github.com/opennetworkinglab/testvectors-runner/pkg/orchestrator/action"
	"github.com/opennetworkinglab/testvectors-runner/pkg/orchestrator/expectation"
	tv "github.com/stratum/testvectors/proto/testvector"
)

var log = logger.NewLogger()

//ProcessTestVector will parse test vector
func ProcessTestVector(tv1 *tv.TestVector) bool {
	result := true
	//log.Infof("Target ID: %s, Target Address: %s\n", target.TargetId, target.Address)
	for _, tc := range tv1.GetTestCases() {
		result = ProcessTestCase(tc) && result
	}
	return result
}

//ProcessTestCase will go through each test case and execute
func ProcessTestCase(tc *tv.TestCase) bool {
	log.Infof("Test Case ID: %s\n", tc.TestCaseId)
	return processActionGroups(tc.GetActionGroups()) && processExpectations(tc.GetExpectations())
}

//processActionGroups processes sequential, parallel or random actions, combines the results and returns true or false.
func processActionGroups(ags []*tv.ActionGroup) bool {
	actionResult := true
	for _, ag := range ags {
		log.Infof("Action Group ID: %s\n", ag.ActionGroupId)
		var currentResult bool
		switch {
		case ag.GetSequentialActionGroup() != nil:
			sag := ag.GetSequentialActionGroup()
			currentResult = action.ProcessSequentialActionGroup(sag)
		case ag.GetParallelActionGroup() != nil:
			pag := ag.GetParallelActionGroup()
			currentResult = action.ProcessParallelActionGroup(pag)
		case ag.GetRandomizedActionGroup() != nil:
			rag := ag.GetRandomizedActionGroup()
			currentResult = action.ProcessRandomizedActionGroup(rag)

		default:
			log.Infof("Empty Action Group\n")
		}
		actionResult = actionResult && currentResult
	}
	return actionResult
}

//processExpectations processes expectations and combines the results to return true or false.
func processExpectations(exps []*tv.Expectation) bool {
	expectationResult := true
	for _, exp := range exps {
		log.Infof("Expectation ID: %s\n", exp.ExpectationId)
		var currentResult bool
		switch {
		case exp.GetConfigExpectation() != nil:
			ce := exp.GetConfigExpectation()
			currentResult = expectation.ProcessConfigExpectation(ce)
		case exp.GetControlPlaneExpectation() != nil:
			cpe := exp.GetControlPlaneExpectation()
			currentResult = expectation.ProcessControlPlaneExpectation(cpe)
		case exp.GetDataPlaneExpectation() != nil:
			dpe := exp.GetDataPlaneExpectation()
			currentResult = expectation.ProcessDataPlaneExpectation(dpe)
		case exp.GetTelemetryExpectation() != nil:
			te := exp.GetTelemetryExpectation()
			currentResult = expectation.ProcessTelemetryExpectation(te)
		default:
			log.Infof("Empty expectation\n")
		}
		expectationResult = expectationResult && currentResult
	}
	return expectationResult
}
