/*
 * Copyright 2019-present Open Networking Foundation
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/golang/protobuf/proto"

	"github.com/opennetworkinglab/testvectors-runner/pkg/logger"
	"github.com/opennetworkinglab/testvectors-runner/pkg/orchestrator"
	"github.com/opennetworkinglab/testvectors-runner/pkg/test"
	"github.com/opennetworkinglab/testvectors-runner/tests"
	tg "github.com/stratum/testvectors/proto/target"
	tv "github.com/stratum/testvectors/proto/testvector"
)

var log = logger.NewLogger()

// main reads test data and utilize testing package to drive the tests. Currently two types of test data are supported.
// One is Test Vectors (see README for more details) and the other is Go function based tests (see examples under tests folder)
// To run with Test Vectors, specify Test Vector files using either tvFiles or tvDir flag, otherwise specify test function
// names using testNames flag. A target file (tgfile) and a port-map file (portMapFile) are mandatory in both cases.
func main() {
	testNames := flag.String("testNames", "", "Names of the tests to run, separated by comma")
	tvFiles := flag.String("tvFiles", "", "Path to the Test Vector files, separated by comma")
	tvDir := flag.String("tvDir", "", "Directory of Test Vector files")
	tgFile := flag.String("tgFile", "", "Path to the Target file")
	portMapFile := flag.String("portMapFile", "tools/bmv2/port-map.json", "Path to the port-map file")
	logDir := flag.String("logDir", "/tmp", "Location to store logs")
	level := flag.String("logLevel", "warn", "Log Level")
	flag.Parse()
	log.SetLogLevel(*level)
	log.SetLogFolder(*logDir)

	// Read Target file
	tgdata, err := ioutil.ReadFile(*tgFile)
	if err != nil {
		log.InvalidFile("Target File: "+*tgFile, err)
	}
	target := &tg.Target{}
	if err = proto.UnmarshalText(string(tgdata), target); err != nil {
		log.InvalidProtoUnmarshal(reflect.TypeOf(target), err)
	}
	log.Infoln("Target: ", target)

	// Read port-map file
	pmdata, err := ioutil.ReadFile(*portMapFile)
	if err != nil {
		log.InvalidFile("Port Map File: "+*portMapFile, err)
	}
	var portmap map[string]string
	if err = json.Unmarshal(pmdata, &portmap); err != nil {
		log.InvalidJSONUnmarshal(reflect.TypeOf(portmap), err)
	}
	log.Infoln("Port Map: ", portmap)

	// Check if we run with Test Vectors or not
	runTV := false
	if *testNames == "" {
		if *tvFiles == "" && *tvDir == "" {
			log.Fatalf("Please specify test names with -testNames or test vector files with -tvFiles or test vector directory with -tvDir")
		} else {
			runTV = true
		}
	}

	// Build test suite
	testSuite := []testing.InternalTest{}
	if !runTV {
		testNameSlice := strings.Split(*testNames, ",")
		stSuite := createTestSuite(testNameSlice, target)
		testSuite = append(testSuite, stSuite...)

	} else {
		// Get a slice of TV files
		var tvFilesSlice []string
		if *tvDir != "" {
			tvFiles, err := ioutil.ReadDir(*tvDir)
			if err != nil {
				log.InvalidDir(*tvDir, err)
			}
			for _, file := range tvFiles {
				if file.IsDir() {
					log.Infof("Ignoring directory %s", file.Name())
					continue
				}
				tvFilesSlice = append(tvFilesSlice, *tvDir+file.Name())
			}
		} else {
			tvFilesSlice = strings.Split(*tvFiles, ",")
		}
		tvSuite := createTVTestSuite(tvFilesSlice, target)
		testSuite = append(testSuite, tvSuite...)
	}

	test.SetUpSuite(target, portmap)
	var match test.Deps
	code := testing.MainStart(match, testSuite, nil, nil).Run()
	test.TearDownSuite()
	os.Exit(code)
}

// createTestSuite creates and returns a slice of InternalTest using a slice of test names
func createTestSuite(testNameSlice []string, target *tg.Target) []testing.InternalTest {
	testSuite := []testing.InternalTest{}
	for _, testName := range testNameSlice {
		// Tests are Go functions
		f := reflect.ValueOf(tests.Test{}).MethodByName(testName)
		if !f.IsValid() {
			log.Fatalf("Not able to find test with name '%s'\nExiting...\n", testName)
		}
		t := testing.InternalTest{
			Name: testName,
			F: func(t *testing.T) {
				test.SetUpTest()
				f.Interface().(func(*testing.T, *tg.Target))(t, target)
				test.TearDownTest()
			},
		}
		testSuite = append(testSuite, t)
	}
	return testSuite
}

// createTVTestSuite creates and returns a slice of InternalTest using a slice of TV files
func createTVTestSuite(tvFilesSlice []string, target *tg.Target) []testing.InternalTest {
	testSuite := []testing.InternalTest{}
	// Read TV files and add them to the test suite
	for _, tvFile := range tvFilesSlice {
		data, err := ioutil.ReadFile(tvFile)
		if err != nil {
			log.InvalidFile("Test Vector File: "+tvFile, err)
		}
		tv := &tv.TestVector{}
		err = proto.UnmarshalText(string(data), tv)
		if err != nil {
			log.InvalidProtoUnmarshal(reflect.TypeOf(tv), err)
		}

		t := testing.InternalTest{
			Name: strings.Replace(filepath.Base(tvFile), ".pb.txt", "", 1),
			F: func(t *testing.T) {
				test.SetUpTest()
				// Process test cases and add them to the test
				for _, tc := range tv.GetTestCases() {
					t.Run(tc.TestCaseId, func(t *testing.T) {
						test.SetUpTestCase(t, target)
						result := orchestrator.ProcessTestCase(tc)
						test.TearDownTestCase(t, target)
						if !result {
							t.Fail()
						}
					})
				}
				test.TearDownTest()
			},
		}
		testSuite = append(testSuite, t)
	}
	return testSuite
}
