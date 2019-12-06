/*
 * Copyright 2019-present Open Networking Foundation
 *
 * SPDX-License-Identifier: Apache-2.0
 */

/*
Package setup implements suite, test and test case setup functions
*/
package setup

import (
	"github.com/opennetworkinglab/testvectors-runner/pkg/framework/dataplane"
	"github.com/opennetworkinglab/testvectors-runner/pkg/framework/gnmi"
	"github.com/opennetworkinglab/testvectors-runner/pkg/framework/p4rt"

	"github.com/opennetworkinglab/testvectors-runner/pkg/logger"
	pm "github.com/stratum/testvectors/proto/portmap"
	tg "github.com/stratum/testvectors/proto/target"
)

var log = logger.NewLogger()

//Suite includes steps for setting up a test suite
func Suite(target *tg.Target, dpMode string, matchType string, portmap *pm.PortMap) {
	log.Info("Setting up test suite...")
	log.Infof("Target: %s", target)
	// Create data plane
	dataplane.CreateDataPlane(dpMode, matchType, portmap)
	gnmi.Init(target)
	p4rt.Init(target, dpMode, portmap)
}

//Test includes steps for setting up a test
func Test() {
	log.Info("Setting up test...")
}

//TestCase includes steps for setting up a test case
func TestCase() {
	log.Info("Setting up test case...")
	// FIXME: only start packet capture if needed
	dataplane.Capture()
}
