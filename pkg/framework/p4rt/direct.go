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

	v1 "github.com/abhilashendurthi/p4runtime/proto/p4/v1"
)

type directPacketIn struct {
	scv streamChannel
}

func (d directPacketIn) ProcessPacketIn(exp *v1.PacketIn) bool {
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
