/*
 * Copyright 2019-present Open Networking Foundation
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package test

import (
	"io"
	"strings"
	"testing"

	"github.com/opennetworkinglab/testvectors-runner/pkg/framework/dataplane"
	"github.com/opennetworkinglab/testvectors-runner/pkg/framework/gnmi"
	"github.com/opennetworkinglab/testvectors-runner/pkg/framework/p4rt"

	"github.com/opennetworkinglab/testvectors-runner/pkg/logger"
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

//SetUpSuite includes steps for setting up test suite
func SetUpSuite(target *tg.Target) {
	log.Infoln("Setting up test suite...")
	log.Infof("Target: %s", target)

	gnmi.Init(target)
	p4rt.Init(target)
}

//TearDownSuite includes steps for tearing down test suite
func TearDownSuite() {
	log.Infoln("Tearing down test suite...")
	gnmi.TearDown()
	p4rt.TearDown()
}

//SetUpTest includes steps for setting up a test
func SetUpTest() {
	log.Infoln("Setting up test...")
}

//TearDownTest includes steps for tearing down a test
func TearDownTest() {
	log.Infoln("Tearing down test...")
}

//SetUpTestCase includes steps for setting up a test case
func SetUpTestCase(t *testing.T, target *tg.Target) {
	log.Debugln("Setting up test case...")
	// FIXME: only start packet capture if needed
	dataplane.Capture()
}

//TearDownTestCase includes steps for tearing down a test case
func TearDownTestCase(t *testing.T, target *tg.Target) {
	log.Debugln("Tearing down test case...")
	dataplane.Stop()
	log.Infoln(strings.Repeat("*", 100))
}
