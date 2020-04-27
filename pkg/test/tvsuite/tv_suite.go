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
	"bytes"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
	"text/template"

	"github.com/golang/protobuf/proto"
	"github.com/stratum/testvectors-runner/pkg/logger"
	"github.com/stratum/testvectors-runner/pkg/orchestrator/testvector"
	"github.com/stratum/testvectors-runner/pkg/test/setup"
	"github.com/stratum/testvectors-runner/pkg/test/teardown"
	tv "github.com/stratum/testvectors/proto/testvector"
)

var log = logger.NewLogger()

//TVSuite struct stores a list of testvector file names, list of template file names and template config file
type TVSuite struct {
	TvFiles        []string
	TemplateFiles  []string
	TemplateConfig string
}

// Create builds and returns a slice of testing.InternalTest from a slice of Test Vector files.
// It iterates through Test Vector files and template files and for each test case it wraps around ProcessTestCase
// to build anonymous functions for testing.InternalTest.
func (tv TVSuite) Create() []testing.InternalTest {
	log.Debug("In Create")
	testSuite := []testing.InternalTest{}
	// Read TV files and add them to the test suite
	for _, tvFile := range tv.TvFiles {
		tv := getTVFromFile(tvFile)
		t := getInternalTest(tvFile, tv)
		testSuite = append(testSuite, t)
	}
	// Read TV template files and add them to the test suite
	for _, templateFile := range tv.TemplateFiles {
		tv := getTVFromTemplateFile(templateFile, tv.TemplateConfig)
		t := getInternalTest(templateFile, tv)
		testSuite = append(testSuite, t)
	}
	return testSuite
}

func getInternalTest(tvFile string, tv *tv.TestVector) testing.InternalTest {
	return testing.InternalTest{
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
}

// getTVFromFile reads Test Vector file with given file name and returns Test Vectors.
func getTVFromFile(fileName string) *tv.TestVector {
	log.Debug("In getTVFromFile")
	tvdata, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalf("Error opening test vector file: %s\n%s", fileName, err)
	}
	testvector := &tv.TestVector{}
	if err = proto.UnmarshalText(string(tvdata), testvector); err != nil {
		log.Fatalf("Error parsing proto message of type %T from file %s\n%s", testvector, fileName, err)
	}
	return testvector
}

// getTVFromTemplateFile reads Template file, config file and returns the converted Test Vectors.
func getTVFromTemplateFile(templateFile string, templateConfigFile string) *tv.TestVector {
	log.Debug("In getTVFromTemplateFile")
	tvdata, err := ioutil.ReadFile(templateFile)
	if err != nil {
		log.Fatalf("Error opening test vector file: %s\n%s", templateFile, err)
	}
	t := template.Must(template.New("tv.tmpl").Parse(string(tvdata)))
	jsondata, err := ioutil.ReadFile(templateConfigFile)
	if err != nil {
		log.Fatalf("Error opening template config file: %s\n%s", templateConfigFile, err)
	}
	m := map[string]interface{}{}
	if err := json.Unmarshal([]byte(jsondata), &m); err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, m); err != nil {
		panic(err)
	}
	testvector := &tv.TestVector{}
	if err = proto.UnmarshalText(buf.String(), testvector); err != nil {
		log.Fatalf("Error parsing proto message of type %T from file %s\n%s", testvector, templateFile, err)
	}
	return testvector
}
