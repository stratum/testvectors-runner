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

//ProcessTestCase will go through each test case and combine the results from processActionGroups and processExpectations to return true or false.
func ProcessTestCase(tc *tv.TestCase) bool {
	log.Infof("Test Case ID: %s\n", tc.TestCaseId)
	return processActionGroups(tc.GetActionGroups()) && processExpectations(tc.GetExpectations())
}

//processActionGroups calls ProcessActionGroup method for each action group in the list, combines all the results to return true or false.
func processActionGroups(ags []*tv.ActionGroup) bool {
	actionResult := true
	for _, ag := range ags {
		log.Infof("Action Group ID: %s\n", ag.ActionGroupId)
		currentResult := action.ProcessActionGroup(ag)
		actionResult = actionResult && currentResult
	}
	return actionResult
}

//processExpectations calls ProcessAExpectation method for each expectation in the list, combines all the results to return true or false.
func processExpectations(exps []*tv.Expectation) bool {
	expectationResult := true
	for _, exp := range exps {
		log.Infof("Expectation ID: %s\n", exp.ExpectationId)
		currentResult := expectation.ProcessExpectation(exp)
		expectationResult = expectationResult && currentResult
	}
	return expectationResult
}
