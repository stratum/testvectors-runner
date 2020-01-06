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

	"github.com/opennetworkinglab/testvectors-runner/pkg/logger"
	tg "github.com/stratum/testvectors/proto/target"
)

var log = logger.NewLogger()

//PktTimeout for receiving all packets
const PktTimeout = 3 * time.Second

var (
	scv      streamChannel
	p4rtConn connection
)

//Init starts a P4Runtime client and runs go routines to send and receive stream channel messages from P4Runtime stream channel client
func Init(target *tg.Target) {
	log.Debug("In p4_oper Init")
	p4rtConn = connect(target)
	scv = getStreamChannel(p4rtConn.client)
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

//ProcessPacketOutOperation sends packet to stream channel client.
func ProcessPacketOutOperation(po *v1.PacketOut) bool {
	var deviceID uint64 = 1
	electionID := &v1.Uint128{High: 1, Low: 5}
	if scv.getMasterArbitrationLock(deviceID, electionID) {
		log.Info("Sending packet")
		log.Debugf("Packet info: %s", po)
		scv.pktOutChan <- po
		return true
	}
	return false
}

//ProcessPacketIn verifies if the packet received is same as expected packet.
func ProcessPacketIn(exp *v1.PacketIn) bool {
	select {
	case ret := <-scv.pktInChan:
		log.Debug("In ProcessPacketIn Case PktInChan")
		return verifyPacketIn(exp, ret)
	case <-time.After(PktTimeout):
		log.Error("Timed out waiting for packet in")
		return false
	}
}
