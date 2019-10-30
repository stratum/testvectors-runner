/*
 * Copyright 2019-present Open Networking Foundation
 *
 * SPDX-License-Identifier: Apache-2.0
 */

/*
Package packet implements basic packet send/receive functions
*/
package packet

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
	"github.com/opennetworkinglab/testvectors-runner/pkg/logger"
)

// Match is used by CheckRawPacket
type Match uint8

// Match values for CheckRawPacket
const (
	Exact   = Match(0x1)
	Contain = Match(0x2)
)

var (
	portMap map[string]string
	// Packet check timeout
	pktCheckTimeout       = 2 * time.Second
	snapshotLen     int32 = 2048
	promiscuous           = false
	// Path for saving pcap files
	pcapPath = "/tmp/"
	// Sleep for 0.5s for each packet verification retry
	packetCheckSleep = 500000000 * time.Nanosecond
	// Maximum duration for packet capturing
	maxTimeout = 1 * time.Hour
	// Map that keeps track of all ongoing captures
	captures sync.Map
	// A wait group which keeps track of all ongoing captures
	wg  sync.WaitGroup
	log = logger.NewLogger()
)

// Init intializes the port map
func Init(portmap map[string]string) {
	portMap = portmap
}

// StartCapture is used to capture the packet and save to a pcap file.
// It takes as arguments the name of the interface for packet captureing and
// a timeout which specifies the duration of the capture.
// When timeout is set to -1*time.Second, it'll use maxTimeout instead.
// It returns a time.Timer which by default is set to timeout and could be
// used to control the duration of the capture.
// Captured packets are saved to a pcap file under pcapPath with the interface
// name as the file name.
// If packet captures on the interface sepcified has already started, it updates
// the timer of the ongoing capture and returns the updated timer
func StartCapture(iface string, timeout time.Duration) *time.Timer {
	if timeout == -1*time.Second {
		timeout = maxTimeout
	}
	// Check if packet capturing on this interface has already started
	if timer, ok := captures.Load(iface); ok {
		// Update the timer and return it
		log.Debugf("Packet capturing already started on %s\n", iface)
		timer.(*time.Timer).Reset(timeout)
		return timer.(*time.Timer)
	}
	// Create pcap file for saving captured packets
	pcapFile := fmt.Sprintf("%s%s.pcap", pcapPath, iface)
	f, _ := os.Create(pcapFile)
	log.Debugf("Saving capture results to %s", pcapFile)
	w := pcapgo.NewWriter(f)
	if err := w.WriteFileHeader(uint32(snapshotLen), layers.LinkTypeEthernet); err != nil {
		log.Fatal(err)
	}

	// Open the device for capturing
	handle, err := pcap.OpenLive(iface, snapshotLen, promiscuous, -1*time.Second)
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
	captures.Store(iface, timer)
	wg.Add(1)
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
			case <-timer.C:
				// Stop capturing on timeout
				log.Debugf("Stop capturing on interface %s...", iface)
				wg.Done()
				return
			}
		}
	}()
	return timer
}

// SendRawPacket is used to send a raw packet to a specific interface.
// It takes as arguments the name of the interface and a slice of byte which
// represents the packet payload.
// It returns true if the packet was successfully sent and false otherwise.
func SendRawPacket(iface string, pkt []byte) bool {
	// Open the device for sending packet
	handle, err := pcap.OpenLive(iface, snapshotLen, promiscuous, -1*time.Second)
	if err != nil {
		log.Errorln(err)
		return false
	}
	defer handle.Close()
	log.Infof("Sending packet to interface %s\n", iface)
	log.Debugf("Packet info: % x\n", pkt)
	if err := handle.WritePacketData(pkt); err != nil {
		log.Errorln(err)
		return false
	}
	return true
}

// CheckRawPacket verifies if packets captured in pcap file are as expected.
// It takes as arguments the name of the interface, a slice of packets with each packet
// represented by a slice of bytes, timeout and match type.
// It verifies that the packets captured in the pcap file match the ones specified in
// pkts within the given timeout. When pkts is empty it verifies that no packet has been
// received.
// When "Exact" is used as match type, it returns true if packets captured are exactly
// the same as pkts including the order. Otherwise it returns false.
// When "Contain" is used as match type, it returns true if packets captured contain pkts.
// Otherwise if returns false.
func CheckRawPacket(iface string, pkts [][]byte, timeout time.Duration, match Match) bool {
	timer := time.After(timeout)
	pcapFile := fmt.Sprintf("%s%s.pcap", pcapPath, iface)
	log.Debugf("Expecting %d packets captured on interface %s\n", len(pkts), iface)
	result := false
	for {
		// Sleep between checks in the loop
	recheck:
		time.Sleep(packetCheckSleep)
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
				log.Errorln(err)
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
					log.Warnf("Expecting %d packets but only captured %d on interface %s...", len(pkts), matchedNum, iface)
					// Recheck until timeout
					result = false
					goto recheck
				} else if err != nil {
					log.Errorln(err)
					// Recheck until timeout
					result = false
					goto recheck
				}
				capturedNum++
				if !bytes.Equal(pkt, packet.Data()) {
					switch match {
					case Exact:
						// Packets don't match, check failed
						log.Errorf("Payloads of packet #%d don't match\n", capturedNum)
						log.Debugf("\nExpected payload: % x\nCaptured payload: % x\n", pkt, packet.Data())
						// No need for recheck in this case
						handle.Close()
						return false
					case Contain:
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
			switch match {
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
			case Contain:
				handle.Close()
				return true
			}
			handle.Close()
		}
	}
}

//StopAllCaptures stops all goroutines
func StopAllCaptures() {
	captures.Range(func(iface, timer interface{}) bool {
		// Stop each capture by resetting its timer to 1 nanosecond
		timer.(*time.Timer).Reset(1 * time.Nanosecond)
		captures.Delete(iface)
		return true
	})
	// Wait for all captures to finish
	wg.Wait()
}

//StartCapturesWithPortMap calls StartCapture for all ports in the port map.
func StartCapturesWithPortMap() bool {
	for _, intf := range portMap {
		log.Debugf("Capturing packets on interface %s\n", intf)
		StartCapture(intf, -1*time.Second)
	}
	return true
}

//SendRawPacketsWithPortMap finds the port in the port map and calls SendRawPacket with the port for each packet.
func SendRawPacketsWithPortMap(pkts [][]byte, port uint32) bool {
	log.Infof("Sending packets to port %d\n", port)
	result := true
	for _, pkt := range pkts {
		result = SendRawPacket(portMap[fmt.Sprint(port)], pkt) && result
	}
	return result
}

//CheckRawPacketsWithPortMap finds the ports in the port map and calls CheckRawPacket for each port.
func CheckRawPacketsWithPortMap(pkts [][]byte, ports []uint32, match Match) bool {
	result := true
	for _, port := range ports {
		log.Infof("Checking packets on port %d\n", port)
		result = CheckRawPacket(portMap[fmt.Sprint(port)], pkts, pktCheckTimeout, match) && result
	}
	return result
}
