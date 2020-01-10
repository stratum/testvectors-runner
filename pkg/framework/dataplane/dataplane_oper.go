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

// CreateDataPlane takes the dataplane mode, packet match type and port-map file name as arguments
// and creates one dataplane instance (only direct dataplane is supported at the moment, loopback
// mode support coming soon) for packet sending/receiving/verification.
func CreateDataPlane(mode string, matchType string, portMap map[string]string) {
	var match Match
	switch matchType {
	case "exact":
		match = Exact
	case "in":
		match = In
	default:
		log.Fatalf("Unknown data plane match type: %s", matchType)
	}
	switch mode {
	case "direct":
		log.Infof("Creating direct data plane with match type: %s and port map: %s\n", matchType, portMap)
		dp = createDirectDataPlane(portMap, match)
	case "loopback":
		log.Infof("Creating loopback data plane with match type: %s and port map: %s\n", matchType, portMap)
		dp = createLoopbackDataPlane(portMap, match)
	default:
		log.Fatalf("Unknown data plane mode: %s", mode)
	}
}

//ProcessTrafficStimulus sends packets to specific ports
func ProcessTrafficStimulus(pkts [][]byte, port uint32) bool {
	log.Debug("In ProcessTrafficStimulus")
	if dp == nil {
		log.Error("data plane does not exist")
		return false
	}
	return dp.send(pkts, port)
}

//ProcessTrafficExpectation verifies that packets arrived at specific ports
func ProcessTrafficExpectation(pkts [][]byte, ports []uint32) bool {
	log.Debug("In ProcessTrafficExpectation")
	if dp == nil {
		log.Error("data plane does not exist")
		return false
	}
	return dp.verify(pkts, ports)
}

//Capture starts packet capturing
func Capture() bool {
	log.Debug("In Capture")
	if dp == nil {
		log.Error("data plane does not exist")
		return false
	}
	return dp.capture()
}

//Stop stops packet capturing
func Stop() bool {
	log.Debug("In Stop")
	if dp == nil {
		log.Error("data plane does not exist")
		return false
	}
	return dp.stop()
}
