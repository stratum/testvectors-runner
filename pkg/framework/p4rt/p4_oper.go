/*
 * Copyright 2019-present Open Networking Foundation
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package p4rt

import (
	"bytes"
	"context"
	"time"

	v1 "github.com/abhilashendurthi/p4runtime/proto/p4/v1"

	"github.com/opennetworkinglab/testvectors-runner/pkg/logger"
	tg "github.com/stratum/testvectors/proto/target"
)

var log = logger.NewLogger()

var (
	//SCV Description
	scv      streamChannel
	p4rtConn connection
)

type connection struct {
	ctx       context.Context
	client    v1.P4RuntimeClient
	connError error
	cancel    context.CancelFunc
}

//StreamChannelVar definition
type streamChannel struct {
	sc                                   v1.P4Runtime_StreamChannelClient
	scError                              error
	cancel                               context.CancelFunc
	masterArbRecvChan, masterArbSendChan chan *v1.MasterArbitrationUpdate
	pktInChan                            chan *v1.PacketIn
	pktOutChan                           chan *v1.PacketOut
	genericStreamMessageChannel          chan *v1.StreamMessageResponse
}

//Init starts a P4Runtime client and runs go routines to send and receive stream channel messages from P4Runtime stream channel client
func Init(target *tg.Target) {
	p4rtConn = connect(target)
	scv = getStreamChannel(p4rtConn.client)
}

//TearDown closes the stream channel client
func TearDown() {
	log.Traceln("In p4_oper tear down")
	scv.cancel()
	if scv.sc != nil {
		err := scv.sc.CloseSend()
		if err != nil {
			log.Warnln("Error closing the stream channel:", err)
		}
	}
	p4rtConn.cancel()
}

//ProcessP4WriteRequest processes the write request
func ProcessP4WriteRequest(wreq *v1.WriteRequest, wres *v1.WriteResponse) bool {
	log.Traceln("In ProcessP4WriteRequest")
	if wreq == nil {
		return false
	}

	lock := getMasterArbitrationLock(scv, wreq.DeviceId, wreq.ElectionId)

	if lock {
		log.Infoln("Sending P4 write request")
		log.Debugf("Write request: %s", wreq)
		ctx := context.Background()
		resp, err := p4rtConn.client.Write(ctx, wreq)
		if err != nil {
			log.Errorf("err:%s\n", err)
			return false
		}
		log.Debugf("resp:%s\n", resp)
		return true
	}
	return false
}

//ProcessP4PipelineConfigOperation processes the forwarding request.
func ProcessP4PipelineConfigOperation(req *v1.SetForwardingPipelineConfigRequest, res *v1.SetForwardingPipelineConfigResponse) bool {
	log.Traceln("In ProcessP4PipelineConfigOperation")
	if req == nil {
		return false
	}
	lock := getMasterArbitrationLock(scv, req.DeviceId, req.ElectionId)
	if lock {
		log.Infoln("Sending P4 pipeline config")
		log.Tracef("Pipeline config: %s", req)
		ctx := context.Background()
		resp, err := p4rtConn.client.SetForwardingPipelineConfig(ctx, req)
		if err != nil {
			log.Errorf("err:%s\n", err)
			return false
		}
		log.Debugf("resp:%s\n", resp)
		return true
	}
	return false
}

//ProcessPacketOutOperation sends packet to stream channel client.
func ProcessPacketOutOperation(po *v1.PacketOut) bool {
	log.Traceln("In ProcessP4 Packet Out")
	var deviceID uint64 = 1
	electionID := &v1.Uint128{High: 1, Low: 5}
	lock := getMasterArbitrationLock(scv, deviceID, electionID)
	if lock {
		log.Infoln("Sending packet")
		log.Debugf("Packet info: %s", po)
		scv.pktOutChan <- po
		return true
	}
	return false
}

//ProcessPacketIn verifies if the packet received is same as expected packet.
func ProcessPacketIn(exp *v1.PacketIn) bool {
	packetMatched := false

	select {
	case ret := <-scv.pktInChan:
		log.Traceln("In ProcessPacketIn Case PktInChan")
		if bytes.Equal(ret.GetPayload(), exp.GetPayload()) {
			packetMatched = true
			log.Infof("Received packet matches")
			log.Debugf("Packet info: %s", ret)
		} else {
			log.Warningf("Packets don't match\nExpected: % x\nActual  : % x\n", exp.GetPayload(), ret.GetPayload())
		}
		return packetMatched
	case <-time.After(3 * time.Second):
		log.Errorln("Timed out")
	}

	return packetMatched
}
