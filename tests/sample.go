/*
*Copyright 2019-present Open Networking Foundation
*
*SPDX-License-Identifier: Apache-2.0
 */

package tests

import (
	"testing"
	"time"

	"github.com/opennetworkinglab/testvectors-runner/pkg/logger"
	"github.com/opennetworkinglab/testvectors-runner/pkg/test"
	tg "github.com/opennetworkinglab/testvectors/proto/target"
	"gotest.tools/assert"
)

//Test struct type
type Test struct{}

var (
	pkt = []byte{0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x0, 0x6, 0x7, 0x8, 0x9, 0xa, 0x8, 0x0, 0x45, 0x0, 0x0, 0x6a, 0x0, 0x1, 0x0, 0x0, 0x40, 0x6, 0xf9, 0x39, 0xc0, 0xa8, 0x0, 0x1, 0xc0, 0xa8, 0x0, 0x2, 0x4, 0xd2, 0x0, 0x50, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x50, 0x2, 0x20, 0x0, 0xe4, 0xe5, 0x0, 0x0, 0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa, 0xb, 0xc, 0xd, 0xe, 0xf, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x3b, 0x3c, 0x3d, 0x3e, 0x3f, 0x40, 0x41}
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

//Test1 is a sample test case
func (st Test) Test1(t *testing.T) {
	// Uses gopacket to send and receive packet locally
	test.VerifyRawPacketLocal("lo", "lo", pkt, 3*time.Second)
}
