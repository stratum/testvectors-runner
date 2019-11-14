/*
 * Copyright 2019-present Open Networking Foundation
 *
 * SPDX-License-Identifier: Apache-2.0
 */

/*
Package dataplane implements packet send/receive functions
*/
package dataplane

import (
	"github.com/opennetworkinglab/testvectors-runner/pkg/logger"
)

var log = logger.NewLogger()

// Match is used by verify
type Match uint8

// Match values for verify
const (
	Exact = Match(0x1)
	In    = Match(0x2)
)

// dataPlane interface implements packet send/receive/verify functions
type dataPlane interface {
	// start packet capturing
	capture() bool
	// send packets to a specific port
	send(pkts [][]byte, port uint32) bool
	// verify packets captured on ports
	verify(pkts [][]byte, ports []uint32) bool
	// stop packet capturing
	stop() bool
}

var dp dataPlane

func CreateDataPlane(mode string, portMap map[string]string, match Match) {
	switch mode {
	case "direct":
		dp = createDirectDataPlane(portMap, match)
	default:
		log.Fatalf("Unknown data plane mode: %s", mode)
	}
}

//ProcessTrafficStimulus sends packets to specific ports
func ProcessTrafficStimulus(pkts [][]byte, port uint32) bool {
	if dp == nil {
		log.Errorln("data plane does not exist")
		return false
	}
	log.Traceln("In ProcessTrafficStimulus")
	return dp.send(pkts, port)
}

//ProcessTrafficExpectation verifies that packets arrived at specific ports
func ProcessTrafficExpectation(pkts [][]byte, ports []uint32) bool {
	if dp == nil {
		log.Errorln("data plane does not exist")
		return false
	}
	log.Traceln("In ProcessTrafficExpectation")
	return dp.verify(pkts, ports)
}

//Capture starts packet capturing
func Capture() bool {
	if dp == nil {
		log.Errorln("data plane does not exist")
		return false
	}
	return dp.capture()
}

//Stop stops packet capturing
func Stop() bool {
	if dp == nil {
		log.Errorln("data plane does not exist")
		return false
	}
	return dp.stop()
}
