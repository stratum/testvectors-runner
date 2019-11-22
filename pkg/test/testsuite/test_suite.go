/*
 * Copyright 2019-present Open Networking Foundation
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package testsuite

import (
	"reflect"
	"testing"

	"github.com/opennetworkinglab/testvectors-runner/pkg/logger"
	"github.com/opennetworkinglab/testvectors-runner/pkg/test/setup"
	"github.com/opennetworkinglab/testvectors-runner/pkg/test/teardown"
	"github.com/opennetworkinglab/testvectors-runner/tests"
)

var log = logger.NewLogger()

type IntTestSuite struct {
	TestNames []string
}

//type Test struct{}

func (ts IntTestSuite) Create() []testing.InternalTest {
	testSuite := []testing.InternalTest{}
	for _, testName := range ts.TestNames {
		// Tests are Go functions
		f := reflect.ValueOf(tests.Test{}).MethodByName(testName)
		if !f.IsValid() {
			log.Fatalf("Not able to find test with name '%s'\nExiting...\n", testName)
		}
		t := testing.InternalTest{
			Name: testName,
			F: func(t *testing.T) {
				setup.Test()
				f.Interface().(func(*testing.T))(t)
				teardown.Test()
			},
		}
		testSuite = append(testSuite, t)
	}
	return testSuite
}