/*
 * Copyright 2019-present Open Networking Foundation
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package dataplane

import (
	pm "github.com/stratum/testvectors/proto/portmap"
	"time"

	v1 "github.com/abhilashendurthi/p4runtime/proto/p4/v1"
	"github.com/opennetworkinglab/testvectors-runner/pkg/framework/p4rt"
	"github.com/opennetworkinglab/testvectors-runner/pkg/utils/common"
)

type loopbackDataPlane struct {
	portmap *pm.PortMap
	match   Match
	// Maximum duration for packet capturing
	maxTimeout time.Duration
}

// createLoopbackDataPlane creates a data plane instance which utilizes packet-out/packet-in to
// mimic data plane packets sending/receiving
func createLoopbackDataPlane(portmap *pm.PortMap, match Match) *loopbackDataPlane {
	ldp := loopbackDataPlane{}
	ldp.portmap = portmap
	ldp.match = match
	ldp.maxTimeout = 1 * time.Hour
	return &ldp
}

// captureOnPort starts packet capturing on all ports specified in portmap.
// In loopback mode p4rt package implements packet capturing functions so there
// is nothing needs to be done here.
// It takes as arguments a timeout which specifies the duration of the capture.
// When timeout is set to -1*time.Second, it'll use maxTimeout instead.
func (ldp *loopbackDataPlane) captureOnPorts(timeout time.Duration) {
}

// sendOnPort sends a raw packet to a specific port via packet-out.
// It takes as arguments the port number and a slice of byte which
// represents the packet payload.
// It returns true if the packet was successfully sent and false otherwise.
func (ldp *loopbackDataPlane) sendOnPort(port uint32, pkt []byte) bool {
	log.Infof("Sending packet to port %d\n", port)
	log.Debugf("Packet info: % x\n", pkt)
	po := convertToPktOut(port, pkt)
	return p4rt.ProcessPacketOutOperation(po)
}

// verifyOnPort verifies if packets captured on sepcific port are as expected.
// It takes as arguments the name of the port and a slice of packets with each packet
// represented by a slice of bytes.
// It verifies that the packets captured on specified port match the ones specified in
// pkts. When pkts is empty it verifies that no packet has been received.
func (ldp *loopbackDataPlane) verifyOnPort(port uint32, pkts [][]byte) bool {
	log.Debugf("Expecting %d packets captured on port %s", len(pkts), port)
	result := true
	for _, pkt := range pkts {
		pi := convertToPktIn(port, pkt)
		result = result && p4rt.ProcessPacketIn(pi)
	}
	// Still need to check for unexpected packets
	//TODO: Check if this PacketIn is valid
	pi := convertToPktIn(port, nil)
	result = result && p4rt.ProcessPacketIn(pi)
	return result
}

//stop stops all captures
func (ldp *loopbackDataPlane) stop() bool {
	// In loopback mode p4rt package implements packet capturing functions so there
	// is nothing needs to be done here.
	return true
}

//capture starts packet capturing
func (ldp *loopbackDataPlane) capture() bool {
	for _, entry := range ldp.portmap.GetEntries() {
		portNumber := entry.GetPortNumber()
		//Commented below section as we start capture on all ports regardless of its type.
		//This if for verifying no packets are captured on ingress port during verifyOnPort()
		/*
			portType := entry.GetPortType()
			if portType == pm.Entry_IN {
				// We don't capture packets on this port
				continue
			}*/
		log.Debugf("Capturing packets on port %d\n", portNumber)
	}
	ldp.captureOnPorts(-1 * time.Second)
	return true
}

//send calls sendOnPort for each packet
func (ldp *loopbackDataPlane) send(pkts [][]byte, port uint32) bool {
	log.Infof("Sending packets to port %d\n", port)
	result := true
	entry := getPortMapEntryByPortNumber(ldp.portmap, port)
	if entry == nil {
		log.Fatalf("Failed to find portmap entry that has port number %d", port)
	}
	portType := entry.GetPortType()
	if portType == pm.Entry_OUT {
		// We shouldn't send packets to this port
		log.Fatalf("Port %d could only be used as egress to switch", port)
	}
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
		entry := getPortMapEntryByPortNumber(ldp.portmap, port)
		if entry != nil {
			//Commented below section in order to verify that no packets are captured on ingress ports.
			//To verify no packets are captured on ingress ports, traffic expectation should have port number and empty packet.
			//verifyOnPort should return true on time out if traffic expectation has empty packet or no packet
			/*portType := entry.GetPortType()
			 if portType == pm.Entry_IN {
				 // We shouldn't capture packets on this port
				 log.Fatalf("Port %d could only be used as ingress to switch", port)
			 }*/
			result = result || ldp.verifyOnPort(port, pkts)
		} else {
			log.Fatalf("Failed to find portmap entry that has port number %d", port)
		}
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
