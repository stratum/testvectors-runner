/*
 * Copyright 2019-present Open Networking Foundation
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package orchestrator

import (
	tv "github.com/stratum/testvectors/proto/testvector"
)

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
	return ProcessActionGroups(tc.GetActionGroups()) && ProcessExpectations(tc.GetExpectations())
}

//ProcessActionGroups processes sequential, parallel or random actions, combines the results and returns true or false.
func ProcessActionGroups(ags []*tv.ActionGroup) bool {
	actionResult := true
	for _, ag := range ags {
		log.Infof("Action Group ID: %s\n", ag.ActionGroupId)
		var currentResult bool
		switch {
		case ag.GetSequentialActionGroup() != nil:
			sag := ag.GetSequentialActionGroup()
			currentResult = ProcessSequentialActionGroup(sag)
		case ag.GetParallelActionGroup() != nil:
			pag := ag.GetParallelActionGroup()
			currentResult = ProcessParallelActionGroup(pag)
		case ag.GetRandomizedActionGroup() != nil:
			rag := ag.GetRandomizedActionGroup()
			currentResult = ProcessRandomizedActionGroup(rag)

		default:
			log.Infof("Empty Action Group\n")
		}
		actionResult = actionResult && currentResult
	}
	return actionResult
}

//ProcessExpectations processes expectations and combines the results to return true or false.
func ProcessExpectations(exps []*tv.Expectation) bool {
	expectationResult := true
	for _, exp := range exps {
		log.Infof("Expectation ID: %s\n", exp.ExpectationId)
		var currentResult bool
		switch {
		case exp.GetConfigExpectation() != nil:
			ce := exp.GetConfigExpectation()
			currentResult = ProcessConfigExpectation(ce)
		case exp.GetControlPlaneExpectation() != nil:
			cpe := exp.GetControlPlaneExpectation()
			currentResult = ProcessControlPlaneExpectation(cpe)
		case exp.GetDataPlaneExpectation() != nil:
			dpe := exp.GetDataPlaneExpectation()
			currentResult = ProcessDataPlaneExpectation(dpe)
		case exp.GetTelemetryExpectation() != nil:
			te := exp.GetTelemetryExpectation()
			currentResult = ProcessTelemetryExpectation(te)
		default:
			log.Infof("Empty expectation\n")
		}
		expectationResult = expectationResult && currentResult
	}
	return expectationResult
}
