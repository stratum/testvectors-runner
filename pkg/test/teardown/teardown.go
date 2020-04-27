/*
 * Copyright 2019-present Open Networking Foundation
 *
 * SPDX-License-Identifier: Apache-2.0
 */

/*
Package teardown implements suite, test and test case teardown functions
*/
package teardown

import (
	"strings"

	"github.com/stratum/testvectors-runner/pkg/framework/dataplane"
	"github.com/stratum/testvectors-runner/pkg/framework/gnmi"
	"github.com/stratum/testvectors-runner/pkg/framework/p4rt"
	"github.com/stratum/testvectors-runner/pkg/logger"
)

var log = logger.NewLogger()

//Suite includes steps for tearing down a test suite
func Suite() {
	log.Info("Tearing down test suite...")
	gnmi.TearDown()
	p4rt.TearDown()
}

//Test includes steps for tearing down a test
func Test() {
	log.Info("Tearing down test...")
}

//TestCase includes steps for tearing down a test case
func TestCase() {
	log.Info("Tearing down test case...")
	dataplane.Stop()
	log.Info(strings.Repeat("*", 100))
}
