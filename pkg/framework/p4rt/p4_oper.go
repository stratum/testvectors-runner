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

	"github.com/stratum/testvectors-runner/pkg/logger"
	"github.com/stratum/testvectors-runner/pkg/utils/common"
	pm "github.com/stratum/testvectors/proto/portmap"
	tg "github.com/stratum/testvectors/proto/target"
)

var log = logger.NewLogger()
var electionID = &v1.Uint128{High: 1, Low: 5}
var deviceID uint64 = 1
var channelSize = 20

//PktTimeout for receiving all packets
const PktTimeout = 3 * time.Second

var (
	s        pktIOInterface
	scv      streamChannel
	p4rtConn connection
)

type pktIOInterface interface {
	ProcessPacketIn(*v1.PacketIn) bool
	ProcessPacketOut(*v1.PacketOut) bool
}

//Init starts a P4Runtime client and runs go routines to send and receive stream channel messages from P4Runtime stream channel client
func Init(target *tg.Target, dpMode string, portmap *pm.PortMap) {
	log.Debug("In p4_oper Init")
	p4rtConn = connect(target)
	scv = getStreamChannel(p4rtConn.client)

	switch dpMode {
	case "direct":
		s = &directPacketIO{scv}
	case "loopback":
		pktChans := make(map[string]chan *v1.PacketIn, channelSize)
		for _, entry := range portmap.GetEntries() {
			portNumber := entry.GetPortNumber()
			pktChans[common.GetStr(portNumber)] = make(chan *v1.PacketIn, channelSize)
		}
		pktChans["generic"] = make(chan *v1.PacketIn, channelSize)
		s = &loopbackPacketIO{scv, pktChans}
		go sort(scv.pktInChan, pktChans)

	default:
		log.Fatalf("Unknown data plane mode: %s", dpMode)
	}
}

//TearDown closes the stream channel client
func TearDown() {
	log.Debug("In p4_oper tear down")
	scv.Close()
	p4rtConn.cancel()
}

//ProcessP4WriteRequest sends the write request to switch
func ProcessP4WriteRequest(wreq *v1.WriteRequest, wres *v1.WriteResponse) bool {
	if wreq == nil {
		return false
	}
	if scv.getMasterArbitrationLock(wreq.DeviceId, wreq.ElectionId) {
		resp := p4rtConn.Write(wreq)
		return verifyWriteResp(wres, resp)
	}
	return false
}

//ProcessP4PipelineConfigOperation sends SetForwardingPipelineConfigRequest to switch
func ProcessP4PipelineConfigOperation(req *v1.SetForwardingPipelineConfigRequest, res *v1.SetForwardingPipelineConfigResponse) bool {
	if req == nil {
		return false
	}
	if scv.getMasterArbitrationLock(req.DeviceId, req.ElectionId) {
		resp := p4rtConn.SetForwardingPipelineConfig(req)
		return verifySetForwardingPipelineConfigResp(res, resp)
	}
	return false
}

//ProcessPacketIn verifies if the packet received is same as expected packet.
func ProcessPacketIn(exp *v1.PacketIn) bool {
	return s.ProcessPacketIn(exp)
}

//ProcessPacketOut sends packet to stream channel client.
func ProcessPacketOut(po *v1.PacketOut) bool {
	return s.ProcessPacketOut(po)
}
