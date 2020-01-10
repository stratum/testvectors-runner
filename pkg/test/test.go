/*
 * Copyright 2019-present Open Networking Foundation
 *
 * SPDX-License-Identifier: Apache-2.0
 */

/*
Package test implements functions to create and run go tests
*/
package test

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/opennetworkinglab/testvectors-runner/pkg/logger"
	"github.com/opennetworkinglab/testvectors-runner/pkg/test/setup"
	"github.com/opennetworkinglab/testvectors-runner/pkg/test/teardown"
	"github.com/opennetworkinglab/testvectors-runner/pkg/test/testsuite"
	"github.com/opennetworkinglab/testvectors-runner/pkg/test/tvsuite"
	tg "github.com/stratum/testvectors/proto/target"
)

var log = logger.NewLogger()

// Deps implements testDeps interface used by MainStart function
type Deps struct{}

func (Deps) ImportPath() string                          { return "" }
func (Deps) MatchString(pat, str string) (bool, error)   { return true, nil }
func (Deps) StartCPUProfile(io.Writer) error             { return nil }
func (Deps) StartTestLog(io.Writer)                      {}
func (Deps) StopCPUProfile()                             {}
func (Deps) StopTestLog() error                          { return nil }
func (Deps) WriteHeapProfile(io.Writer) error            { return nil }
func (Deps) WriteProfileTo(string, io.Writer, int) error { return nil }

//Suite interface defines Create method for converting tv files, test names to go test type
type Suite interface {
	Create() []testing.InternalTest
}

//CreateSuite returns a slice of InternalTest based on go test names or testvector directory name.
func CreateSuite(testNames string, tvDir string, tvName string) []testing.InternalTest {
	var testSuite Suite
	switch {
	case testNames != "":
		var ts testsuite.IntTestSuite
		ts.TestNames = strings.Split(testNames, ",")
		testSuite = ts
	case tvDir != "":
		var tvs tvsuite.TVSuite
		re, _ := regexp.Compile("^" + tvName + "\\.pb.txt$")
		log.Debugf("Test Vectors file regex: %s", re)
		tvs.TvFiles = getFiles(tvDir, re)
		log.Debugf("Test Vectors to run: %s", tvs.TvFiles)
		testSuite = tvs
	}
	return testSuite.Create()
}

//getFiles walks through given directory and returns list of all files in the directory
func getFiles(tvDir string, re *regexp.Regexp) []string {
	log.Debug("In getFiles")
	var tvFilesSlice []string
	err := filepath.Walk(tvDir,
		func(filePath string, fileInfo os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if re.MatchString(fileInfo.Name()) {
				tvFilesSlice = append(tvFilesSlice, filePath)
			}
			return nil
		})
	if err != nil {
		log.Fatalf("Error opening directory: %s\n%s", tvDir, err)
	}
	return tvFilesSlice
}

//getTarget reads the given file and converts it to target proto.
//panics if file is invalid
func getTarget(fileName string) *tg.Target {
	// Read Target file
	tgdata, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalf("Error opening target file: %s\n%s", fileName, err)
	}
	target := &tg.Target{}
	if err = proto.UnmarshalText(string(tgdata), target); err != nil {
		log.Fatalf("Error parsing proto message of type %T from file %s\n%s", target, fileName, err)
	}
	log.Info("Target: ", target)
	return target
}

//Run calls suite setup, teardown and runs all tests in the testSuite against given target
func Run(tgFile string, dpMode string, matchType string, portMapFile string, testSuite []testing.InternalTest) {
	log.Debug("In Run")
	target := getTarget(tgFile)
	portMap := getPortMap(portMapFile)
	setup.Suite(target, dpMode, matchType, portMap)
	var match Deps
	code := testing.MainStart(match, testSuite, nil, nil).Run()
	teardown.Suite()
	os.Exit(code)
}

//getPortMap reads the given file and converts it to portMap.
//panics if file is invalid
func getPortMap(fileName string) map[string]string {
	// Read port-map file
	pmdata, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalf("Error opening port map file: %s\n%s", fileName, err)
	}
	var portMap map[string]string
	if err = json.Unmarshal(pmdata, &portMap); err != nil {
		log.Fatalf("Error parsing json data of type %T from file %s\n%s", portMap, fileName, err)
	}
	return portMap
}
