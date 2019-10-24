/*
*Copyright 2019-present Open Networking Foundation
*
*SPDX-License-Identifier: Apache-2.0
 */

package orchestrator

import (
	tg "github.com/stratum/testvectors/proto/target"
	tv "github.com/stratum/testvectors/proto/testvector"
)

//Target global variable
var Target *tg.Target

//ProcessTestVector will parse test vector
func ProcessTestVector(tv1 *tv.TestVector, tg1 *tg.Target) bool {
	if Target = tg1; Target == nil {
		log.Errorln("Unspecified target, cannot execute tests")
		return false
	}
	result := true
	log.Infof("Target ID: %s, Target Address: %s\n", Target.TargetId, Target.Address)
	for _, tc := range tv1.GetTestCases() {
		result = ProcessTestCase(tc, Target) && result
	}
	return result
}

//ProcessTestCase will go through each test case and execute
func ProcessTestCase(tc *tv.TestCase, tg *tg.Target) bool {
	if Target = tg; Target == nil {
		log.Errorln("Unspecified target, cannot execute tests")
		return false
	}
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
