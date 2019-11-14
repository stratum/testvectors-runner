/*
 * Copyright 2019-present Open Networking Foundation
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package tests

import (
	"testing"

	"github.com/opennetworkinglab/testvectors-runner/pkg/logger"
	tg "github.com/stratum/testvectors/proto/target"
	"gotest.tools/assert"
)

//Test struct type
type Test struct{}

var (
	log = logger.NewLogger()
)

//TestCase1 is a sample sub test
func TestCase1(t *testing.T) {
	log.Infoln("Starting TestCase1")
	assert.Equal(t, true, false)
}

//TestCase2 is a sample sub test
func TestCase2(t *testing.T) {
	log.Infoln("Starting TestCase2")
	t.Error("Fail this test")
	t.Fail()
}

//Test0 is a sample test case
func (st Test) Test0(t *testing.T, target *tg.Target) {
	log.Infoln(target.Address)
	t.Run("Test Case 1", TestCase1)
	t.Run("Test Case 2", TestCase2)
}
