/*
 * Copyright 2019-present Open Networking Foundation
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package test

import (
	"io"
	"io/ioutil"
	"os"
	"reflect"
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

//TestSuite description
type Suite interface {
	Create() []testing.InternalTest
}

//CreateSuite description
func CreateSuite(testNames string, tvDir string, tvFiles string) []testing.InternalTest {
	var testSuite Suite
	switch {
	case testNames != "":
		var ts testsuite.IntTestSuite
		ts.TestNames = strings.Split(testNames, ",")
		testSuite = ts
	case tvDir != "":
		var tvs tvsuite.TVSuite
		tvs.TvFiles = getFiles(tvDir)
		testSuite = tvs
	case tvFiles != "":
		var tvs tvsuite.TVSuite
		tvs.TvFiles = strings.Split(tvFiles, ",")
		testSuite = tvs
	}
	return testSuite.Create()
}

func getFiles(tvDir string) []string {
	var tvFilesSlice []string
	tvFiles, err := ioutil.ReadDir(tvDir)
	if err != nil {
		log.InvalidDir(tvDir, err)
	}
	for _, file := range tvFiles {
		if file.IsDir() {
			log.Infof("Ignoring directory %s", file.Name())
			continue
		}
		tvFilesSlice = append(tvFilesSlice, tvDir+file.Name())
	}
	return tvFilesSlice
}

func getTarget(fileName string) *tg.Target {
	// Read Target file
	tgdata, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.InvalidFile("Target File: "+fileName, err)
	}
	target := &tg.Target{}
	if err = proto.UnmarshalText(string(tgdata), target); err != nil {
		log.InvalidProtoUnmarshal(reflect.TypeOf(target), err)
	}
	log.Infoln("Target: ", target)
	return target
}

//Run description
func Run(tgFile string, testSuite []testing.InternalTest) {
	target := getTarget(tgFile)
	setup.Suite(target)
	var match Deps
	code := testing.MainStart(match, testSuite, nil, nil).Run()
	teardown.Suite()
	os.Exit(code)
}
