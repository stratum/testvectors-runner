/*
 * Copyright 2019-present Open Networking Foundation
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package dataplane

import (
	"time"

	pm "github.com/stratum/testvectors/proto/portmap"

	v1 "github.com/p4lang/p4runtime/go/p4/v1"
	"github.com/stratum/testvectors-runner/pkg/framework/p4rt"
	"github.com/stratum/testvectors-runner/pkg/utils/common"
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

// sendOnPort sends a raw packet to a specific port via packet-out.
// It takes as arguments the port number and a slice of byte which
// represents the packet payload.
// It returns true if the packet was successfully sent and false otherwise.
func (ldp *loopbackDataPlane) sendOnPort(port uint32, pkt []byte) bool {
	log.Infof("Sending packet to port %d\n", port)
	log.Debugf("Packet info: % x\n", pkt)
	po := convertToPktOut(port, pkt)
	return p4rt.ProcessPacketOut(po)
}

// verifyOnPort verifies if packets captured on sepcific port are as expected.
// It takes as arguments the name of the port and a slice of packets with each packet
// represented by a slice of bytes.
// It verifies that the packets captured on specified port match the ones specified in
// pkts. When pkts is empty it verifies that no packet has been received.
func (ldp *loopbackDataPlane) verifyOnPort(port uint32, pkts [][]byte) bool {
	log.Debugf("Expecting %d packets captured on port %d", len(pkts), port)
	result := true
	for _, pkt := range pkts {
		pi := convertToPktIn(port, pkt)
		result = result && p4rt.ProcessPacketIn(pi)
	}
	// Unexpected packets can be checked with empty packet in dataplane expectation
	// pi := convertToPktIn(port, nil)
	// result = result && p4rt.ProcessPacketIn(pi)
	return result
}

//stop stops all captures
func (ldp *loopbackDataPlane) stop() bool {
	//This function is empty because packet capture in loopback is done as part of p4rt packet-ins but stop() still has to be defined as part of the interface definition
	return true
}

//capture starts packet capturing
func (ldp *loopbackDataPlane) capture() bool {
	//This function is empty because packet capture in loopback is done as part of p4rt packet-ins but capture() still has to be defined as part of the interface definition
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
	//MetadataId 1 represents egress_physical_port, 2 represents cpuLoopbackMode, 3 represents padding of packet_out controller_packet_metadata based on the P4 program
	//cpuLoopbackMode == 2 represents a dataplane packet meant to be sent to egress
	//TODO: make the metadata id configurable to allow loopback for any P4 program
	po.Metadata = []*v1.PacketMetadata{
		{MetadataId: 1, Value: common.GetByteSlice(port, 2)},
		{MetadataId: 2, Value: common.GetByteSlice(2, 1)},
		{MetadataId: 3, Value: common.GetByteSlice(0, 1)},
		{MetadataId: 4, Value: common.GetByteSlice(48897, 2)},
	}
	log.Debugf("Packet info: %s", po)
	return po
}

func convertToPktIn(port uint32, pkt []byte) *v1.PacketIn {
	pi := &v1.PacketIn{}
	pi.Payload = pkt
	pi.Metadata = []*v1.PacketMetadata{
		//MetadataId 1 represents ingress_physical_port of packet_in controller_packet_metadata based on the P4 program
		//TODO: make the metadata id configurable to allow loopback for any P4 program
		{MetadataId: 1, Value: common.GetByteSlice(port, 2)},
	}
	return pi
}
