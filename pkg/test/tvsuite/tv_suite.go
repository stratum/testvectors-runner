/*
 * Copyright 2019-present Open Networking Foundation
 *
 * SPDX-License-Identifier: Apache-2.0
 */

/*
Package tvsuite implements Create function to convert testvector files to go tests
*/
package tvsuite

import (
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/opennetworkinglab/testvectors-runner/pkg/logger"
	"github.com/opennetworkinglab/testvectors-runner/pkg/orchestrator/testvector"
	"github.com/opennetworkinglab/testvectors-runner/pkg/test/setup"
	"github.com/opennetworkinglab/testvectors-runner/pkg/test/teardown"
	tv "github.com/stratum/testvectors/proto/testvector"
)

var log = logger.NewLogger()

//TVSuite struct stores a list of testvector file names
type TVSuite struct {
	TvFiles []string
}

// Create builds and returns a slice of testing.InternalTest from a slice of Test Vector files.
// It iterates through Test Vector files and for each test case it wraps around ProcessTestCase
// to build anonymous functions for testing.InternalTest.
func (tv TVSuite) Create() []testing.InternalTest {
	testSuite := []testing.InternalTest{}
	// Read TV files and add them to the test suite
	for _, tvFile := range tv.TvFiles {

		tv := getTVFromFile(tvFile)
		t := testing.InternalTest{
			Name: strings.Replace(filepath.Base(tvFile), ".pb.txt", "", 1),
			F: func(t *testing.T) {
				setup.Test()
				// Process test cases and add them to the test
				for _, tc := range tv.GetTestCases() {
					t.Run(tc.TestCaseId, func(t *testing.T) {
						setup.TestCase()
						result := testvector.ProcessTestCase(tc)
						teardown.TestCase()
						if !result {
							t.Fail()
						}
					})
				}
				teardown.Test()
			},
		}
		testSuite = append(testSuite, t)
	}
	return testSuite
}

// getTVFromFile reads Test Vector file with given file name and returns Test Vectors.
func getTVFromFile(fileName string) *tv.TestVector {
	tvdata, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.InvalidFile("Test Vector File: "+fileName, err)
	}
	testvector := &tv.TestVector{}
	if err = proto.UnmarshalText(string(tvdata), testvector); err != nil {
		log.InvalidProtoUnmarshal(reflect.TypeOf(testvector), err)
	}
	return testvector
}
