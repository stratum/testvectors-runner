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
	"encoding/json"
	"io/ioutil"
	"reflect"

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

func CreateDataPlane(mode string, matchType string, portMapFile string) {
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
		// Read port-map file
		pmdata, err := ioutil.ReadFile(portMapFile)
		if err != nil {
			log.InvalidFile("Port Map File: "+portMapFile, err)
		}
		var portMap map[string]string
		if err = json.Unmarshal(pmdata, &portMap); err != nil {
			log.InvalidJSONUnmarshal(reflect.TypeOf(portMap), err)
		}
		log.Infof("Creating direct data plane with match type: %s and port map: %s\n", matchType, portMap)
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
