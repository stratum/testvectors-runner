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
	"github.com/stratum/testvectors-runner/pkg/utils/common"
)

type loopbackPacketIn struct {
	scv      streamChannel
	pktChans map[string]chan *v1.PacketIn
}

func (l loopbackPacketIn) ProcessPacketIn(exp *v1.PacketIn) bool {
	//FIXME: instead of first Metadata object, find metadata that matches with ingress port id
	ingressPort := common.GetStr(exp.GetMetadata()[0].GetValue())
	if _, ok := l.pktChans[ingressPort]; !ok {
		ingressPort = "generic"
	}
	select {
	case ret := <-l.pktChans[ingressPort]:
		log.Debug("In ProcessPacketIn Case PktInChan")
		return verifyPacketIn(exp, ret)
	case <-time.After(PktTimeout):
		if exp.GetPayload() == nil {
			return true
		}
		log.Error("Timed out waiting for packet in")
		return false
	}
}

func sort(pktInChan chan *v1.PacketIn, pktChans map[string]chan *v1.PacketIn) {
	for {
		packet := <-pktInChan
		log.Debugf("Caught packet in sort %v", packet)
		//FIXME: instead of first Metadata object, find metadata that matches with ingress port id
		ingressPort := common.GetStr(packet.GetMetadata()[0].GetValue())
		if val, ok := pktChans[ingressPort]; ok {
			log.Debugf("Added packet to channel with port %s", ingressPort)
			val <- packet
		} else {
			//pktChans[value] = make(chan *v1.PacketIn)
			log.Debugf("Added packet to channel with port generic")
			pktChans["generic"] <- packet
		}
	}
}
