/*
 * Copyright 2019-present Open Networking Foundation
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package dataplane

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/pcapgo"
)

type directDataPlane struct {
	portMap map[string]string
	match   Match
	// Packet check timeout
	pktCheckTimeout time.Duration
	snapshotLen     int32
	promiscuous     bool
	// Path for saving pcap files
	pcapPath string
	// Sleep for 0.5s for each packet verification retry
	packetCheckSleep time.Duration
	// Maximum duration for packet capturing
	maxTimeout time.Duration
	// Map that keeps track of all ongoing captures
	captures sync.Map
	// A wait group which keeps track of all ongoing captures
	wg sync.WaitGroup
}

// createDirectDataPlane creates a data plane instance which utilizes gopacket to
// send/receive packets directly to physical interfaces on the host.
func createDirectDataPlane(portMap map[string]string, match Match) *directDataPlane {
	ddp := directDataPlane{}
	ddp.portMap = portMap
	ddp.match = match
	ddp.pktCheckTimeout = 2 * time.Second
	ddp.snapshotLen = 2048
	ddp.promiscuous = false
	ddp.pcapPath = "/tmp/"
	ddp.packetCheckSleep = 500000000 * time.Nanosecond
	ddp.maxTimeout = 1 * time.Hour
	return &ddp
}

// captureOnInterface is used to capture the packet and save to a pcap file.
// It takes as arguments the name of the interface for packet captureing and
// a timeout which specifies the duration of the capture.
// When timeout is set to -1*time.Second, it'll use maxTimeout instead.
// It returns a time.Timer which by default is set to timeout and could be
// used to control the duration of the capture.
// Captured packets are saved to a pcap file under pcapPath with the interface
// name as the file name.
// If packet captures on the interface sepcified has already started, it updates
// the timer of the ongoing capture and returns the updated timer
func (ddp *directDataPlane) captureOnInterface(iface string, timeout time.Duration) *time.Timer {
	if timeout == -1*time.Second {
		timeout = ddp.maxTimeout
	}
	// Check if packet capturing on this interface has already started
	if timer, ok := ddp.captures.Load(iface); ok {
		// Update the timer and return it
		log.Debugf("Packet capturing already started on %s\n", iface)
		timer.(*time.Timer).Reset(timeout)
		return timer.(*time.Timer)
	}
	// Create pcap file for saving captured packets
	pcapFile := fmt.Sprintf("%s%s.pcap", ddp.pcapPath, iface)
	f, _ := os.Create(pcapFile)
	log.Debugf("Saving capture results to %s", pcapFile)
	w := pcapgo.NewWriter(f)
	if err := w.WriteFileHeader(uint32(ddp.snapshotLen), layers.LinkTypeEthernet); err != nil {
		log.Fatal(err)
	}

	// Open the device for capturing
	handle, err := pcap.OpenLive(iface, ddp.snapshotLen, ddp.promiscuous, -1*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	// Only capture packets with "in" direction and ingore packets sent by ourselves
	if err := handle.SetDirection(pcap.DirectionIn); err != nil {
		log.Fatal(err)
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	timer := time.NewTimer(timeout)
	// Keep track of the capture using a global map
	ddp.captures.Store(iface, timer)
	ddp.wg.Add(1)
	// Start the packet capturing loop in a goroutine and return the timer
	go func() {
		defer f.Close()
		defer handle.Close()
		for {
			select {
			case packet := <-packetSource.Packets():
				// Save captured packets to file
				log.Infof("Caught packet on interface %s\n", iface)
				log.Debugf("Packet info: %s\n", packet)
				if err := w.WritePacket(packet.Metadata().CaptureInfo, packet.Data()); err != nil {
					log.Fatal(err)
				}
				// TODO: also save packets to a buffer for later verification
			case <-timer.C:
				// Stop capturing on timeout
				log.Debugf("Stop capturing on interface %s...", iface)
				ddp.wg.Done()
				return
			}
		}
	}()
	return timer
}

// sendOnInterface is used to send a raw packet to a specific interface.
// It takes as arguments the name of the interface and a slice of byte which
// represents the packet payload.
// It returns true if the packet was successfully sent and false otherwise.
func (ddp *directDataPlane) sendOnInterface(iface string, pkt []byte) bool {
	// Open the device for sending packet
	handle, err := pcap.OpenLive(iface, ddp.snapshotLen, ddp.promiscuous, -1*time.Second)
	if err != nil {
		log.Error(err)
		return false
	}
	defer handle.Close()
	log.Infof("Sending packet to interface %s\n", iface)
	log.Debugf("Packet info: % x\n", pkt)
	if err := handle.WritePacketData(pkt); err != nil {
		log.Error(err)
		return false
	}
	return true
}

// verifyOnInterface verifies if packets captured in pcap file are as expected.
// It takes as arguments the name of the interface and a slice of packets with each packet
// represented by a slice of bytes.
// It verifies that the packets captured in the pcap file match the ones specified in
// pkts within the timeout. When pkts is empty it verifies that no packet has been
// received.
// When "Exact" is used as match type, it returns true if packets captured are exactly
// the same as pkts including the order. Otherwise it returns false.
// When "In" is used as match type, it returns true if packets captured contain pkts.
// Otherwise if returns false.
func (ddp *directDataPlane) verifyOnInterface(iface string, pkts [][]byte) bool {
	timer := time.After(ddp.pktCheckTimeout)
	// TODO: read packets from buffer instead of file
	// Also see TODO in captureOnInterface()
	pcapFile := fmt.Sprintf("%s%s.pcap", ddp.pcapPath, iface)
	log.Debugf("Expecting %d packets captured on interface %s\n", len(pkts), iface)
	result := false
	for {
		// Sleep between checks in the loop
	recheck:
		time.Sleep(ddp.packetCheckSleep)
		select {
		case <-timer:
			log.Debugf("Timeout exceeded, stop checking packet on interface %s...", iface)
			if result {
				log.Infof("Packet check passed on interface %s...", iface)
			} else {
				log.Errorf("Packet check failed on interface %s...", iface)
			}
			return result
		default:
			result = true
			// Open pcap file
			handle, err := pcap.OpenOffline(pcapFile)
			if err != nil {
				log.Error(err)
				return false
			}
			log.Debugf("Reading packet from %s...\n", pcapFile)
			packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
			capturedNum := 0
			matchedNum := 0
			// Loop through expected packets
			for _, pkt := range pkts {
				// Get next captured packet
			nextpacket:
				packet, err := packetSource.NextPacket()
				if err == io.EOF {
					// We've reached the end of the pcap file but we expect more packets captured
					log.Debugf("Expecting %d packets but only captured %d on interface %s...", len(pkts), matchedNum, iface)
					// Recheck until timeout
					result = false
					goto recheck
				} else if err != nil {
					log.Error(err)
					// Recheck until timeout
					result = false
					goto recheck
				}
				capturedNum++
				if !bytes.Equal(pkt, packet.Data()) {
					switch ddp.match {
					case Exact:
						// Packets don't match, check failed
						log.Errorf("Payloads of packet #%d don't match\n", capturedNum)
						log.Debugf("\nExpected payload: % x\nCaptured payload: % x\n", pkt, packet.Data())
						// No need for recheck in this case
						handle.Close()
						return false
					case In:
						// Packets don't match, ignore it
						log.Debugf("Ingoring unmached packet: % x\n", packet.Data())
						goto nextpacket
					}
				} else {
					// Packets match
					log.Debugf("Payloads of packet #%d match\n", capturedNum)
					matchedNum++
				}
			}
			switch ddp.match {
			case Exact:
				// Check if there are more captured packets than expected
				for packet := range packetSource.Packets() {
					log.Debugf("Unexpected packet on interface %s: %s\n", iface, packet)
					capturedNum++
				}
				if capturedNum > len(pkts) {
					log.Errorf("Expecting %d packets but captured %d on interface %s...", len(pkts), capturedNum, iface)
					// No need for recheck in this case
					handle.Close()
					return false
				}
			case In:
				handle.Close()
				return true
			}
			handle.Close()
		}
	}
}

//stop stops all goroutines
func (ddp *directDataPlane) stop() bool {
	ddp.captures.Range(func(iface, timer interface{}) bool {
		// Stop each capture by resetting its timer to 1 nanosecond
		timer.(*time.Timer).Reset(1 * time.Nanosecond)
		ddp.captures.Delete(iface)
		return true
	})
	// Wait for all captures to finish
	ddp.wg.Wait()
	return true
}

//capture calls captureOnInterface for all ports in the port map.
func (ddp *directDataPlane) capture() bool {
	for _, intf := range ddp.portMap {
		log.Debugf("Capturing packets on interface %s\n", intf)
		ddp.captureOnInterface(intf, -1*time.Second)
	}
	return true
}

//send finds the port in the port map and calls sendOnInterface with the port for each packet.
func (ddp *directDataPlane) send(pkts [][]byte, port uint32) bool {
	log.Infof("Sending packets to port %d\n", port)
	result := true
	for _, pkt := range pkts {
		result = ddp.sendOnInterface(ddp.portMap[fmt.Sprint(port)], pkt) && result
	}
	return result
}

//verify finds the ports in the port map and calls verifyOnInterface for each port.
func (ddp *directDataPlane) verify(pkts [][]byte, ports []uint32) bool {
	result := true
	for _, port := range ports {
		log.Infof("Checking packets on port %d\n", port)
		result = ddp.verifyOnInterface(ddp.portMap[fmt.Sprint(port)], pkts) && result
	}
	return result
}
