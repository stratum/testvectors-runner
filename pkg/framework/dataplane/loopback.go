/*
 * Copyright 2019-present Open Networking Foundation
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package dataplane

import (
	//"bytes"

	"time"

	v1 "github.com/abhilashendurthi/p4runtime/proto/p4/v1"
	"github.com/opennetworkinglab/testvectors-runner/pkg/framework/p4rt"
	"github.com/opennetworkinglab/testvectors-runner/pkg/utils/common"
)

type loopbackDataPlane struct {
	portMap map[string]string
	match   Match
	// Packet check timeout
	pktCheckTimeout time.Duration
	// Maximum duration for packet capturing
	maxTimeout time.Duration
	// Timer that controls the duration of packet capturing
	captureTimer *time.Timer
	// TODO: channels for buffering captured packets
}

// createLoopbackDataPlane creates a data plane instance which utilizes packet-out/packet-in to
// mimic data plane packets sending/receiving
func createLoopbackDataPlane(portMap map[string]string, match Match) *loopbackDataPlane {
	ldp := loopbackDataPlane{}
	ldp.portMap = portMap
	ldp.match = match
	ldp.pktCheckTimeout = 2 * time.Second
	ldp.maxTimeout = 1 * time.Hour
	return &ldp
}

// captureOnPort starts packet capturing on all ports specified in portMap.
// It saves packet-ins to a channel for future processing.
// It takes as arguments a timeout which specifies the duration of the capture.
// When timeout is set to -1*time.Second, it'll use maxTimeout instead.
func (ldp *loopbackDataPlane) captureOnPorts(timeout time.Duration) {
	if timeout == -1*time.Second {
		timeout = ldp.maxTimeout
	}
	ldp.captureTimer = time.NewTimer(timeout)
	// Start the packet capturing loop in a goroutine
	go func() {
		for {
			select {
			// TODO: replace the following commented block: get packet-ins from some channel
			/*
				case packet := <-packetSource.Packets():
					log.Infof("Caught packet on port %s\n", port)
					log.Debugf("Packet info: %s\n", packet)
					// Save packet to channel, use different channels for different ports
			*/
			case <-ldp.captureTimer.C:
				// Stop capturing on timeout
				log.Debugf("Stop packet capturing")
				return
			// The default case is only for passing CI and will be updated/removed once implemetation completes
			default:
				time.Sleep(1 * time.Second)
			}
		}
	}()
}

// sendOnPort sends a raw packet to a specific port via packet-out.
// It takes as arguments the port number and a slice of byte which
// represents the packet payload.
// It returns true if the packet was successfully sent and false otherwise.
func (ldp *loopbackDataPlane) sendOnPort(port uint32, pkt []byte) bool {
	log.Infof("Sending packet to port %d\n", port)
	log.Debugf("Packet info: % x\n", pkt)
	// TODO: send packet to some channel that sends it as a packet-out
	po := convertToPktOut(port, pkt)
	return p4rt.ProcessPacketOutOperation(po)
	//return true
}

// verifyOnPort verifies if packets captured on sepcific port are as expected.
// It takes as arguments the name of the port and a slice of packets with each packet
// represented by a slice of bytes.
// It verifies that the packets captured on specified port match the ones specified in
// pkts within the timeout. When pkts is empty it verifies that no packet has been
// received.
// When "Exact" is used as match type, it returns true if packets captured are exactly
// the same as pkts including the order. Otherwise it returns false.
// When "In" is used as match type, it returns true if packets captured contain pkts.
// Otherwise if returns false.
func (ldp *loopbackDataPlane) verifyOnPort(port uint32, pkts [][]byte) bool {
	timer := time.After(ldp.pktCheckTimeout)
	log.Debugf("Expecting %d packets captured on port %d\n", len(pkts), port)
	result := false
	for {
		select {
		// TODO: read from buffer and compare packets
		/*
			if !bytes.Equal(pkt, packet) {
			}
		*/
		case <-timer:
			log.Debugf("Timeout exceeded, stop checking packet on port %d...", port)
			if result {
				log.Infof("Packet check passed on port %d...", port)
			} else {
				log.Errorf("Packet check failed on port %d...", port)
			}
			return result
		// The default case is only for passing CI and will be updated/removed once implemetation completes
		default:
			if len(pkts) > 0 {
				result := true
				for _, pkt := range pkts {
					pi := convertToPktIn(port, pkt)
					result = result && p4rt.ProcessPacketIn(pi)
				}
				return result
			}
			//TODO: Check if this PacketIn is valid
			pi := convertToPktIn(port, nil)
			return p4rt.ProcessPacketIn(pi)
			//time.Sleep(1 * time.Second)
		}
	}
}

//stop stops all captures
func (ldp *loopbackDataPlane) stop() bool {
	for _, port := range ldp.portMap {
		log.Debugf("Stop packet capturing on port %s\n", port)
		// TODO: delete the ACL rules installed for packet capturing
	}
	// Stop packet capturing by resetting the capture timer
	ldp.captureTimer.Reset(1 * time.Nanosecond)
	return true
}

//capture starts packet capturing
func (ldp *loopbackDataPlane) capture() bool {
	for _, port := range ldp.portMap {
		log.Debugf("Capturing packets on port %s\n", port)
		// TODO: install ACL rules to redirect packets received to CPU
	}
	ldp.captureOnPorts(-1 * time.Second)
	return true
}

//send calls sendOnPort for each packet
func (ldp *loopbackDataPlane) send(pkts [][]byte, port uint32) bool {
	log.Infof("Sending packets to port %d\n", port)
	result := true
	for _, pkt := range pkts {
		result = ldp.sendOnPort(port, pkt) && result
	}
	return result
}

//verify calls verifyOnPort for each port
func (ldp *loopbackDataPlane) verify(pkts [][]byte, ports []uint32) bool {
	result := false
	for _, port := range ports {
		log.Infof("Checking packets on port %d\n", port)
		result = result || ldp.verifyOnPort(port, pkts)
	}
	return result
}

func convertToPktOut(port uint32, pkt []byte) *v1.PacketOut {
	po := &v1.PacketOut{}
	po.Payload = pkt
	po.Metadata = []*v1.PacketMetadata{{MetadataId: 1, Value: common.GetUint32(port)}}
	return po
}

func convertToPktIn(port uint32, pkt []byte) *v1.PacketIn {
	pi := &v1.PacketIn{}
	pi.Payload = pkt
	pi.Metadata = []*v1.PacketMetadata{
		{MetadataId: 1, Value: common.GetUint32(port)},
		//{MetadataId: 4, Value: common.GetUint32(port)},
	}
	return pi
}
