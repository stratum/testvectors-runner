/*
 * Copyright 2019-present Open Networking Foundation
 *
 * SPDX-License-Identifier: Apache-2.0
 */

/*
Package p4rt implements p4runtime functions
*/
package p4rt

import (
	"time"

	v1 "github.com/p4lang/p4runtime/go/p4/v1"
)

type directPacketIO struct {
	scv streamChannel
}

func (d directPacketIO) ProcessPacketIn(exp *v1.PacketIn) bool {
	select {
	case ret := <-d.scv.pktInChan:
		log.Debug("In ProcessPacketIn Case PktInChan")
		return verifyPacketIn(exp, ret)
	case <-time.After(PktTimeout):
		if exp == nil || exp.GetPayload() == nil {
			return true
		}
		log.Error("Timed out waiting for packet in")
		return false
	}
}

func (d directPacketIO) ProcessPacketOut(po *v1.PacketOut) bool {
	if d.scv.getMasterArbitrationLock(deviceID, electionID) {
		log.Info("Sending packet")
		log.Debugf("Packet info: %s", po)
		scv.pktOutChan <- po
		return true
	}
	return false
}
