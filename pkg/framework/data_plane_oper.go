/*
*Copyright 2019-present Open Networking Foundation
*
*SPDX-License-Identifier: Apache-2.0
 */

package framework

import (
	"github.com/opennetworkinglab/testvectors-runner/pkg/packet"
	tg "github.com/opennetworkinglab/testvectors/proto/target"
)

//ProcessTrafficStimulus sends packets to specific ports
func ProcessTrafficStimulus(target *tg.Target, pkts [][]byte, port uint32) bool {
	log.Traceln("In ProcessTrafficStimulus")
	return packet.SendRawPacketsWithPortMap(pkts, port)
}

//ProcessTrafficExpectation verifies that packets arrived at specific ports
func ProcessTrafficExpectation(target *tg.Target, pkts [][]byte, ports []uint32) bool {
	log.Traceln("In ProcessTrafficExpectation")
	return packet.CheckRawPacketsWithPortMap(pkts, ports, packet.Exact)
}
