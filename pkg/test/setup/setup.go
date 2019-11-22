/*
 * Copyright 2019-present Open Networking Foundation
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package setup

import (
	"github.com/opennetworkinglab/testvectors-runner/pkg/framework/dataplane"
	"github.com/opennetworkinglab/testvectors-runner/pkg/framework/gnmi"
	"github.com/opennetworkinglab/testvectors-runner/pkg/framework/p4rt"

	"github.com/opennetworkinglab/testvectors-runner/pkg/logger"
	tg "github.com/stratum/testvectors/proto/target"
)

var log = logger.NewLogger()

//Suite includes steps for setting up test suite
func Suite(target *tg.Target) {
	log.Infoln("Setting up test suite...")
	log.Infof("Target: %s", target)

	gnmi.Init(target)
	p4rt.Init(target)
}

//Test includes steps for setting up a test
func Test() {
	log.Infoln("Setting up test...")
}

//TestCase includes steps for setting up a test case
func TestCase() {
	log.Debugln("Setting up test case...")
	// FIXME: only start packet capture if needed
	dataplane.Capture()
}
