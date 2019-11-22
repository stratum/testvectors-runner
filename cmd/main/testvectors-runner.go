/*
 * Copyright 2019-present Open Networking Foundation
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"flag"

	"github.com/opennetworkinglab/testvectors-runner/pkg/framework/dataplane"
	"github.com/opennetworkinglab/testvectors-runner/pkg/logger"
	"github.com/opennetworkinglab/testvectors-runner/pkg/test"
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
	dpMode := flag.String("dpMode", "direct", "Data plane mode: 'direct' or 'loopback'")
	matchType := flag.String("matchType", "exact", "Data plane match type: 'exact' or 'in'")
	logDir := flag.String("logDir", "/tmp", "Location to store logs")
	level := flag.String("logLevel", "warn", "Log Level")
	flag.Parse()

	setupLog(*logDir, *level)

	// Create data plane
	dataplane.CreateDataPlane(*dpMode, *matchType, *portMapFile)

	testSuiteSlice := test.CreateSuite(*testNames, *tvDir, *tvFiles)
	test.Run(*tgFile, testSuiteSlice)
}

func setupLog(logDir string, logLevel string) {
	log.SetLogLevel(logLevel)
	log.SetLogFolder(logDir)
}
